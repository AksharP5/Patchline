package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/opencode"
	"github.com/AksharP5/Patchline/internal/snapshot"
)

func snapshotCommand(opts CommonOptions, stdout io.Writer, stderr io.Writer) int {
	result, err := opencode.Discover(opts.ProjectRoot, opts.GlobalConfig, []string(opts.LocalDirs))
	if err != nil {
		fmt.Fprintf(stderr, "failed to discover plugins: %v\n", err)
		return 1
	}

	snapshotDir, candidates := snapshot.ResolveDir(opts.SnapshotDir)
	if snapshotDir == "" {
		fmt.Fprintln(stderr, "snapshot directory not found")
		return 1
	}
	if opts.SnapshotDir == "" && len(candidates) > 0 {
		fmt.Fprintf(stderr, "Using snapshot directory: %s\n", snapshotDir)
	}

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

	store := snapshot.Store{Directory: snapshotDir}
	saved := 0
	localCount := 0
	for _, spec := range result.Plugins {
		if spec.Source == opencode.SourceLocal {
			localCount++
			continue
		}
		installed := "missing"
		if entry, ok := installedByName[spec.Name]; ok {
			installed = entry.Version
		}
		err := store.Save(snapshot.Entry{
			PluginName:        spec.Name,
			PreviousSpec:      spec.DeclaredSpec,
			PreviousInstalled: installed,
			Source:            string(spec.Source),
			Reason:            "snapshot",
			ConfigPath:        spec.ConfigPath,
		})
		if err != nil {
			fmt.Fprintf(stderr, "failed to save snapshot for %s: %v\n", spec.Name, err)
			return 1
		}
		saved++
	}

	if saved == 0 {
		fmt.Fprintln(stdout, "No npm plugins found to snapshot.")
		return 0
	}

	fmt.Fprintf(stdout, "Saved %d snapshot(s) to %s.\n", saved, snapshotDir)
	if localCount > 0 {
		fmt.Fprintln(stdout, "Local plugins are unmanaged and were skipped.")
	}
	return 0
}

func rollbackCommand(opts CommonOptions, pluginName string, stdout io.Writer, stderr io.Writer) int {
	snapshotDir, candidates := snapshot.ResolveDir(opts.SnapshotDir)
	if snapshotDir == "" {
		fmt.Fprintln(stderr, "snapshot directory not found")
		return 1
	}
	if opts.SnapshotDir == "" && len(candidates) > 0 {
		fmt.Fprintf(stderr, "Using snapshot directory: %s\n", snapshotDir)
	}

	store := snapshot.Store{Directory: snapshotDir}
	entry, err := store.Latest(pluginName)
	if err != nil {
		if errors.Is(err, snapshot.ErrSnapshotNotFound) {
			fmt.Fprintf(stderr, "no snapshot found for %s\n", pluginName)
			return 1
		}
		fmt.Fprintf(stderr, "failed to load snapshot: %v\n", err)
		return 1
	}
	if entry.ConfigPath == "" {
		fmt.Fprintln(stderr, "snapshot missing config path")
		return 1
	}

	if err := opencode.UpdatePluginSpec(entry.ConfigPath, pluginName, entry.PreviousSpec); err != nil {
		fmt.Fprintf(stderr, "failed to update config: %v\n", err)
		return 1
	}

	ctx := context.Background()
	cacheDir, cacheCandidates := cache.ResolveDir(opts.CacheDir)
	if cacheDir != "" {
		if _, err := cache.Invalidate(ctx, cacheDir, pluginName); err != nil {
			fmt.Fprintf(stderr, "failed to invalidate cache: %v\n", err)
			return 1
		}
	} else if opts.CacheDir != "" {
		fmt.Fprintf(stderr, "cache directory not found: %s\n", opts.CacheDir)
	} else if len(cacheCandidates) > 0 {
		fmt.Fprintf(stderr, "cache directory not found. Checked: %s\n", strings.Join(cacheCandidates, ", "))
	}

	fmt.Fprintf(stdout, "Restored %s to %s. Run OpenCode to reinstall.\n", pluginName, entry.PreviousSpec)
	return 0
}
