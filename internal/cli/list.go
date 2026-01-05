package cli

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/model"
	"github.com/AksharP5/Patchline/internal/opencode"
)

func listCommand(opts CommonOptions, stdout io.Writer, stderr io.Writer) int {
	result, err := opencode.Discover(opts.ProjectRoot, opts.GlobalConfig, []string(opts.LocalDirs))
	if err != nil {
		fmt.Fprintf(stderr, "failed to discover plugins: %v\n", err)
		return 1
	}

	cacheDir, candidates := cache.ResolveDir(opts.CacheDir)
	cacheEntries := []cache.Entry{}
	ctx := context.Background()
	if cacheDir != "" {
		cacheEntries, err = cache.Detect(ctx, cacheDir)
		if err != nil {
			fmt.Fprintf(stderr, "failed to scan cache directory: %v\n", err)
			return 1
		}
	} else if opts.CacheDir != "" {
		fmt.Fprintf(stderr, "cache directory not found: %s\n", opts.CacheDir)
	} else if len(candidates) > 0 {
		fmt.Fprintf(stderr, "cache directory not found. Checked: %s\n", strings.Join(candidates, ", "))
	}

	plugins := buildPluginList(result.Plugins, cacheEntries)
	renderPluginTable(stdout, plugins)
	printListHints(stdout, plugins, cacheDir == "")
	return 0
}

func buildPluginList(specs []opencode.PluginSpec, cacheEntries []cache.Entry) []model.Plugin {
	cacheByName := make(map[string]cache.Entry, len(cacheEntries))
	for _, entry := range cacheEntries {
		cacheByName[entry.Name] = entry
	}

	plugins := make([]model.Plugin, 0, len(specs))
	for _, spec := range specs {
		plugin := model.Plugin{
			Name:         spec.Name,
			DeclaredSpec: spec.DeclaredSpec,
			Source:       string(spec.Source),
			ConfigPath:   spec.ConfigPath,
			CachePath:    "",
		}

		if spec.Source == opencode.SourceLocal {
			plugin.Status = model.StatusUnmanaged
			plugin.Installed = "local"
			plugin.LocalDirectory = spec.LocalPath
			plugins = append(plugins, plugin)
			continue
		}

		entry, ok := cacheByName[spec.Name]
		if ok {
			plugin.Installed = entry.Version
			plugin.CachePath = entry.Path
			if spec.Pinned != "" && entry.Version != spec.Pinned {
				plugin.Status = model.StatusMismatch
			} else {
				plugin.Status = model.StatusOK
			}
		} else {
			plugin.Installed = "missing"
			plugin.Status = model.StatusMissing
		}

		plugins = append(plugins, plugin)
	}

	sort.Slice(plugins, func(i, j int) bool {
		if plugins[i].Name == plugins[j].Name {
			return plugins[i].Source < plugins[j].Source
		}
		return plugins[i].Name < plugins[j].Name
	})

	return plugins
}

func renderPluginTable(w io.Writer, plugins []model.Plugin) {
	if len(plugins) == 0 {
		fmt.Fprintln(w, "No plugins found.")
		return
	}

	headers := []string{"NAME", "DECLARED", "INSTALLED", "STATUS", "SOURCE"}
	rows := make([][]string, 0, len(plugins))
	for _, plugin := range plugins {
		rows = append(rows, []string{
			plugin.Name,
			plugin.DeclaredSpec,
			plugin.Installed,
			string(plugin.Status),
			plugin.Source,
		})
	}

	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s\n",
		widths[0], headers[0],
		widths[1], headers[1],
		widths[2], headers[2],
		widths[3], headers[3],
		widths[4], headers[4],
	)
	for _, row := range rows {
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s\n",
			widths[0], row[0],
			widths[1], row[1],
			widths[2], row[2],
			widths[3], row[3],
			widths[4], row[4],
		)
	}
}

func printListHints(w io.Writer, plugins []model.Plugin, cacheMissing bool) {
	needsSync := false
	for _, plugin := range plugins {
		if plugin.Status == model.StatusMissing || plugin.Status == model.StatusMismatch {
			needsSync = true
			break
		}
	}

	if needsSync {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Next step: run `patchline sync` to refresh the cache.")
	}

	if cacheMissing {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Tip: pass --cache-dir to point at the OpenCode cache.")
	}
}
