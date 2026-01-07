package cli

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/npm"
	"github.com/AksharP5/Patchline/internal/opencode"
	"github.com/AksharP5/Patchline/internal/snapshot"
)

type upgradeTarget struct {
	Name       string
	ConfigPath string
	Source     string
	Declared   string
	Pinned     string
}

func upgradeCommand(opts CommonOptions, name string, target string, mode string, all bool, stdout io.Writer, stderr io.Writer) int {
	result, err := opencode.Discover(opts.ProjectRoot, opts.GlobalConfig, []string(opts.LocalDirs))
	if err != nil {
		fmt.Fprintf(stderr, "failed to discover plugins: %v\n", err)
		return 1
	}

	targets := selectUpgradeTargets(result.Plugins, name, all)
	if len(targets) == 0 {
		if all {
			fmt.Fprintln(stderr, "no npm plugins found to upgrade")
		} else {
			fmt.Fprintf(stderr, "plugin not found: %s\n", name)
		}
		return 1
	}

	if target == "" && opts.Offline {
		fmt.Fprintln(stderr, "upgrade requires registry access; provide --to or disable --offline")
		return 2
	}

	snapshotDir, candidates := snapshot.ResolveDir(opts.SnapshotDir)
	if snapshotDir == "" {
		fmt.Fprintln(stderr, "snapshot directory not found")
		return 1
	}
	if opts.SnapshotDir == "" && len(candidates) > 0 {
		fmt.Fprintf(stderr, "Using snapshot directory: %s\n", snapshotDir)
	}
	store := snapshot.Store{Directory: snapshotDir}

	ctx := context.Background()
	cacheDir, cacheCandidates := cache.ResolveDir(opts.CacheDir)
	cacheEntries := []cache.Entry{}
	if cacheDir != "" {
		cacheEntries, err = cache.Detect(ctx, cacheDir)
		if err != nil {
			fmt.Fprintf(stderr, "failed to scan cache directory: %v\n", err)
			return 1
		}
	} else if opts.CacheDir != "" {
		fmt.Fprintf(stderr, "cache directory not found: %s\n", opts.CacheDir)
	} else if len(cacheCandidates) > 0 {
		fmt.Fprintf(stderr, "cache directory not found. Checked: %s\n", strings.Join(cacheCandidates, ", "))
	}

	installedByName := map[string]cache.Entry{}
	for _, entry := range cacheEntries {
		installedByName[entry.Name] = entry
	}

	infoCache := map[string]npm.PackageInfo{}
	updated := 0
	skipped := 0
	for _, targetSpec := range targets {
		installedVersion := ""
		installedLabel := "missing"
		if entry, ok := installedByName[targetSpec.Name]; ok {
			installedVersion = entry.Version
			installedLabel = entry.Version
		}
		base := chooseBaseVersion(targetSpec.Pinned, installedVersion)

		resolved := target
		if resolved == "" {
			info, ok := infoCache[targetSpec.Name]
			if !ok {
				info, err = npm.FetchPackageInfo(ctx, targetSpec.Name)
				if err != nil {
					fmt.Fprintf(stderr, "failed to fetch %s: %v\n", targetSpec.Name, err)
					return 1
				}
				infoCache[targetSpec.Name] = info
			}
			resolved, err = npm.SelectTargetVersion(info.Latest, info.Versions, base, npm.UpgradeMode(mode))
			if err != nil {
				fmt.Fprintf(stderr, "failed to resolve target for %s: %v\n", targetSpec.Name, err)
				return 1
			}
		}

		newSpec := fmt.Sprintf("%s@%s", targetSpec.Name, resolved)
		if strings.TrimSpace(targetSpec.Declared) == newSpec {
			fmt.Fprintf(stdout, "Already on %s\n", newSpec)
			skipped++
			continue
		}

		err := store.Save(snapshot.Entry{
			PluginName:        targetSpec.Name,
			PreviousSpec:      targetSpec.Declared,
			PreviousInstalled: installedLabel,
			Source:            targetSpec.Source,
			Reason:            "upgrade",
			ConfigPath:        targetSpec.ConfigPath,
		})
		if err != nil {
			fmt.Fprintf(stderr, "failed to save snapshot for %s: %v\n", targetSpec.Name, err)
			return 1
		}

		if err := opencode.UpdatePluginSpec(targetSpec.ConfigPath, targetSpec.Name, newSpec); err != nil {
			fmt.Fprintf(stderr, "failed to update config for %s: %v\n", targetSpec.Name, err)
			return 1
		}

		if cacheDir != "" {
			if _, err := cache.Invalidate(ctx, cacheDir, targetSpec.Name); err != nil {
				fmt.Fprintf(stderr, "failed to invalidate cache for %s: %v\n", targetSpec.Name, err)
				return 1
			}
		}

		fmt.Fprintf(stdout, "Upgraded %s -> %s\n", targetSpec.Name, newSpec)
		updated++
	}

	if updated == 0 {
		if skipped > 0 {
			fmt.Fprintln(stdout, "All plugins already match the target versions.")
			return 0
		}
		fmt.Fprintln(stdout, "No plugins upgraded.")
		return 0
	}

	fmt.Fprintln(stdout, "")
	fmt.Fprintf(stdout, "Updated %d plugin(s). Run OpenCode to reinstall.\n", updated)
	if skipped > 0 {
		fmt.Fprintf(stdout, "%d plugin(s) already matched the target.\n", skipped)
	}
	return 0
}

func selectUpgradeTargets(specs []opencode.PluginSpec, name string, all bool) []upgradeTarget {
	byName := map[string][]opencode.PluginSpec{}
	for _, spec := range specs {
		if spec.Source == opencode.SourceLocal || spec.ConfigPath == "" {
			continue
		}
		byName[spec.Name] = append(byName[spec.Name], spec)
	}

	if all {
		names := make([]string, 0, len(byName))
		for pluginName := range byName {
			names = append(names, pluginName)
		}
		sort.Strings(names)

		var targets []upgradeTarget
		for _, pluginName := range names {
			targets = append(targets, selectPreferredTargets(byName[pluginName])...)
		}
		return targets
	}

	return selectPreferredTargets(byName[name])
}

func selectPreferredTargets(specs []opencode.PluginSpec) []upgradeTarget {
	if len(specs) == 0 {
		return nil
	}

	order := []opencode.Source{
		opencode.SourceCustom,
		opencode.SourceCustomDir,
		opencode.SourceProject,
		opencode.SourceGlobal,
	}
	for _, source := range order {
		filtered := filterTargets(specs, source)
		if len(filtered) > 0 {
			return filtered
		}
	}
	return nil
}

func filterTargets(specs []opencode.PluginSpec, source opencode.Source) []upgradeTarget {
	seen := map[string]struct{}{}
	out := []upgradeTarget{}
	for _, spec := range specs {
		if spec.Source != source {
			continue
		}
		if spec.ConfigPath == "" {
			continue
		}
		if _, ok := seen[spec.ConfigPath]; ok {
			continue
		}
		seen[spec.ConfigPath] = struct{}{}
		out = append(out, upgradeTarget{
			Name:       spec.Name,
			ConfigPath: spec.ConfigPath,
			Source:     string(spec.Source),
			Declared:   spec.DeclaredSpec,
			Pinned:     spec.Pinned,
		})
	}
	return out
}

func chooseBaseVersion(pinned string, installed string) string {
	if pinned != "" {
		if _, ok := npm.CompareSemver(pinned, pinned); ok {
			return pinned
		}
	}
	return installed
}
