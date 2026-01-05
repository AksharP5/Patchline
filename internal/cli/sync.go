package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/model"
	"github.com/AksharP5/Patchline/internal/opencode"
)

type syncRow struct {
	Name      string
	Declared  string
	Installed string
	Status    string
	Action    string
	Source    string
}

type syncPlan struct {
	Rows            []syncRow
	RefreshTargets  []string
	SkippedUnpinned int
	LocalCount      int
}

func syncCommand(opts CommonOptions, stdout io.Writer, stderr io.Writer) int {
	result, err := opencode.Discover(opts.ProjectRoot, opts.GlobalConfig, []string(opts.LocalDirs))
	if err != nil {
		fmt.Fprintf(stderr, "failed to discover plugins: %v\n", err)
		return 1
	}

	cacheDir, candidates := cache.ResolveDir(opts.CacheDir)
	if cacheDir == "" {
		if opts.CacheDir != "" {
			fmt.Fprintf(stderr, "cache directory not found: %s\n", opts.CacheDir)
		} else if len(candidates) > 0 {
			fmt.Fprintf(stderr, "cache directory not found. Checked: %v\n", candidates)
		} else {
			fmt.Fprintln(stderr, "cache directory not found")
		}
		return 1
	}

	entries, err := cache.Detect(cacheDir)
	if err != nil {
		fmt.Fprintf(stderr, "failed to scan cache directory: %v\n", err)
		return 1
	}

	plan := buildSyncPlan(result.Plugins, entries)
	renderSyncTable(stdout, plan.Rows)

	if len(plan.RefreshTargets) == 0 {
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Cache already matches pinned config.")
		if plan.SkippedUnpinned > 0 {
			fmt.Fprintf(stdout, "Skipped %d unpinned plugin(s).\n", plan.SkippedUnpinned)
		}
		if plan.LocalCount > 0 {
			fmt.Fprintln(stdout, "Local plugins are unmanaged and were skipped.")
		}
		return 0
	}

	refreshed := 0
	missing := 0
	for _, name := range plan.RefreshTargets {
		removed, err := cache.Invalidate(cacheDir, name)
		if err != nil {
			fmt.Fprintf(stderr, "failed to refresh %s: %v\n", name, err)
			return 1
		}
		if len(removed) == 0 {
			missing++
			continue
		}
		refreshed++
	}

	fmt.Fprintln(stdout, "")
	if refreshed > 0 {
		fmt.Fprintf(stdout, "Refreshed %d plugin(s). Run OpenCode to reinstall.\n", refreshed)
	}
	if missing > 0 {
		fmt.Fprintf(stdout, "No cache entry found for %d plugin(s); OpenCode will install on next run.\n", missing)
	}
	if plan.SkippedUnpinned > 0 {
		fmt.Fprintf(stdout, "Skipped %d unpinned plugin(s).\n", plan.SkippedUnpinned)
	}
	if plan.LocalCount > 0 {
		fmt.Fprintln(stdout, "Local plugins are unmanaged and were skipped.")
	}
	return 0
}

func buildSyncPlan(specs []opencode.PluginSpec, entries []cache.Entry) syncPlan {
	installedByName := map[string]cache.Entry{}
	for _, entry := range entries {
		installedByName[entry.Name] = entry
	}

	rows := []syncRow{}
	targets := []string{}
	skippedUnpinned := 0
	localCount := 0

	for _, spec := range specs {
		if spec.Source == opencode.SourceLocal {
			localCount++
			rows = append(rows, syncRow{
				Name:      spec.Name,
				Declared:  spec.DeclaredSpec,
				Installed: "local",
				Status:    string(model.StatusUnmanaged),
				Action:    "skip",
				Source:    string(spec.Source),
			})
			continue
		}

		entry, ok := installedByName[spec.Name]
		installed := "missing"
		if ok {
			installed = entry.Version
		}

		status := string(model.StatusOK)
		action := "noop"
		if spec.Pinned == "" {
			status = string(model.StatusUnknown)
			action = "skip"
			skippedUnpinned++
		} else if installed == "missing" {
			status = string(model.StatusMissing)
			action = "refresh"
			targets = append(targets, spec.Name)
		} else if spec.Pinned != installed {
			status = string(model.StatusMismatch)
			action = "refresh"
			targets = append(targets, spec.Name)
		}

		rows = append(rows, syncRow{
			Name:      spec.Name,
			Declared:  spec.DeclaredSpec,
			Installed: installed,
			Status:    status,
			Action:    action,
			Source:    string(spec.Source),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Name == rows[j].Name {
			return rows[i].Source < rows[j].Source
		}
		return rows[i].Name < rows[j].Name
	})

	return syncPlan{
		Rows:            rows,
		RefreshTargets:  uniqueNames(targets),
		SkippedUnpinned: skippedUnpinned,
		LocalCount:      localCount,
	}
}

func renderSyncTable(w io.Writer, rows []syncRow) {
	if len(rows) == 0 {
		fmt.Fprintln(w, "No plugins found.")
		return
	}

	headers := []string{"NAME", "DECLARED", "INSTALLED", "STATUS", "ACTION", "SOURCE"}
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		values := []string{row.Name, row.Declared, row.Installed, row.Status, row.Action, row.Source}
		for i, value := range values {
			if len(value) > widths[i] {
				widths[i] = len(value)
			}
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
		widths[0], headers[0],
		widths[1], headers[1],
		widths[2], headers[2],
		widths[3], headers[3],
		widths[4], headers[4],
		widths[5], headers[5],
	)

	for _, row := range rows {
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
			widths[0], row.Name,
			widths[1], row.Declared,
			widths[2], row.Installed,
			widths[3], row.Status,
			widths[4], row.Action,
			widths[5], row.Source,
		)
	}
}

func uniqueNames(values []string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
