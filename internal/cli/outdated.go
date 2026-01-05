package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/model"
	"github.com/AksharP5/Patchline/internal/npm"
	"github.com/AksharP5/Patchline/internal/opencode"
)

type outdatedRow struct {
	Name      string
	Declared  string
	Installed string
	Latest    string
	Status    string
	Source    string
}

func outdatedCommand(opts CommonOptions, stdout io.Writer, stderr io.Writer) int {
	if opts.Offline {
		fmt.Fprintln(stderr, "outdated requires registry access; remove --offline")
		return 2
	}

	result, err := opencode.Discover(opts.ProjectRoot, opts.GlobalConfig, []string(opts.LocalDirs))
	if err != nil {
		fmt.Fprintf(stderr, "failed to discover plugins: %v\n", err)
		return 1
	}

	cacheDir, _ := cache.ResolveDir(opts.CacheDir)
	cacheEntries := []cache.Entry{}
	if cacheDir != "" {
		cacheEntries, err = cache.Detect(cacheDir)
		if err != nil {
			fmt.Fprintf(stderr, "failed to scan cache directory: %v\n", err)
			return 1
		}
	}

	installedByName := map[string]cache.Entry{}
	for _, entry := range cacheEntries {
		installedByName[entry.Name] = entry
	}

	latestByName := map[string]string{}
	localCount := 0
	for _, spec := range result.Plugins {
		if spec.Source == opencode.SourceLocal {
			localCount++
			continue
		}
		if _, ok := latestByName[spec.Name]; ok {
			continue
		}
		info, err := npm.FetchPackageInfo(spec.Name)
		if err != nil {
			fmt.Fprintf(stderr, "failed to fetch %s: %v\n", spec.Name, err)
			continue
		}
		latestByName[spec.Name] = info.Latest
	}

	rows := make([]outdatedRow, 0, len(result.Plugins))
	for _, spec := range result.Plugins {
		if spec.Source == opencode.SourceLocal {
			continue
		}

		entry, ok := installedByName[spec.Name]
		installed := "missing"
		if ok {
			installed = entry.Version
		}

		latest := latestByName[spec.Name]
		status := string(model.StatusOK)
		if installed == "missing" {
			status = string(model.StatusMissing)
		} else if latest != "" {
			if cmp, ok := npm.CompareSemver(installed, latest); ok && cmp < 0 {
				status = string(model.StatusOutdated)
			}
		}

		rows = append(rows, outdatedRow{
			Name:      spec.Name,
			Declared:  spec.DeclaredSpec,
			Installed: installed,
			Latest:    latest,
			Status:    status,
			Source:    string(spec.Source),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Name == rows[j].Name {
			return rows[i].Source < rows[j].Source
		}
		return rows[i].Name < rows[j].Name
	})

	renderOutdatedTable(stdout, rows)
	if localCount > 0 {
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Note: local plugins are unmanaged and excluded from outdated checks.")
	}
	return 0
}

func renderOutdatedTable(w io.Writer, rows []outdatedRow) {
	if len(rows) == 0 {
		fmt.Fprintln(w, "No npm plugins found.")
		return
	}

	headers := []string{"NAME", "DECLARED", "INSTALLED", "LATEST", "STATUS", "SOURCE"}
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		values := []string{row.Name, row.Declared, row.Installed, row.Latest, row.Status, row.Source}
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
			widths[3], row.Latest,
			widths[4], row.Status,
			widths[5], row.Source,
		)
	}
}
