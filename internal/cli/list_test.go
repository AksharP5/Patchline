package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/model"
	"github.com/AksharP5/Patchline/internal/opencode"
)

func TestBuildPluginListStatuses(t *testing.T) {
	specs := []opencode.PluginSpec{
		{
			Name:         "local",
			DeclaredSpec: "/tmp/local.js",
			Source:       opencode.SourceLocal,
			LocalPath:    "/tmp/local.js",
		},
		{
			Name:         "missing",
			DeclaredSpec: "missing@1.0.0",
			Pinned:       "1.0.0",
			Source:       opencode.SourceProject,
			ConfigPath:   "/tmp/opencode.json",
		},
		{
			Name:         "mismatch",
			DeclaredSpec: "mismatch@1.0.0",
			Pinned:       "1.0.0",
			Source:       opencode.SourceProject,
			ConfigPath:   "/tmp/opencode.json",
		},
		{
			Name:         "ok",
			DeclaredSpec: "ok@1.0.0",
			Pinned:       "1.0.0",
			Source:       opencode.SourceProject,
			ConfigPath:   "/tmp/opencode.json",
		},
	}

	cacheEntries := []cache.Entry{
		{Name: "mismatch", Version: "2.0.0", Path: "/cache/mismatch"},
		{Name: "ok", Version: "1.0.0", Path: "/cache/ok"},
	}

	plugins := buildPluginList(specs, cacheEntries)
	byName := map[string]model.Plugin{}
	for _, plugin := range plugins {
		byName[plugin.Name] = plugin
	}

	local := byName["local"]
	if local.Status != model.StatusUnmanaged || local.Installed != "local" {
		t.Fatalf("expected local plugin to be unmanaged/local, got status=%s installed=%s", local.Status, local.Installed)
	}

	missing := byName["missing"]
	if missing.Status != model.StatusMissing || missing.Installed != "missing" {
		t.Fatalf("expected missing plugin, got status=%s installed=%s", missing.Status, missing.Installed)
	}

	mismatch := byName["mismatch"]
	if mismatch.Status != model.StatusMismatch || mismatch.Installed != "2.0.0" {
		t.Fatalf("expected mismatch plugin, got status=%s installed=%s", mismatch.Status, mismatch.Installed)
	}

	ok := byName["ok"]
	if ok.Status != model.StatusOK || ok.Installed != "1.0.0" {
		t.Fatalf("expected ok plugin, got status=%s installed=%s", ok.Status, ok.Installed)
	}
}

func TestPrintListHints(t *testing.T) {
	plugins := []model.Plugin{
		{Name: "missing", Status: model.StatusMissing},
	}

	var out bytes.Buffer
	printListHints(&out, plugins, false)
	output := out.String()
	if !strings.Contains(output, "Next step: run `patchline sync`") {
		t.Fatalf("expected sync hint, got %q", output)
	}

	out.Reset()
	printListHints(&out, nil, true)
	output = out.String()
	if !strings.Contains(output, "Tip: pass --cache-dir") {
		t.Fatalf("expected cache-dir hint, got %q", output)
	}
}

func TestListCommandNoPlugins(t *testing.T) {
	root := t.TempDir()
	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    root,
	}

	var out bytes.Buffer
	var err bytes.Buffer
	code := listCommand(opts, &out, &err)
	if code != 0 {
		t.Fatalf("expected success, got %d", code)
	}
	if err.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", err.String())
	}
	if !strings.Contains(out.String(), "No plugins found.") {
		t.Fatalf("expected no plugins message, got %q", out.String())
	}
}
