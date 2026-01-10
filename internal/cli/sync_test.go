package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AksharP5/Patchline/internal/cache"
	"github.com/AksharP5/Patchline/internal/model"
	"github.com/AksharP5/Patchline/internal/opencode"
)

func TestBuildSyncPlan(t *testing.T) {
	specs := []opencode.PluginSpec{
		{Name: "alpha", DeclaredSpec: "alpha@1.0.0", Pinned: "1.0.0", Source: opencode.SourceProject},
		{Name: "beta", DeclaredSpec: "beta", Pinned: "", Source: opencode.SourceProject},
		{Name: "local-plugin", DeclaredSpec: "/tmp/local.js", Source: opencode.SourceLocal},
	}
	entries := []cache.Entry{
		{Name: "alpha", Version: "0.9.0"},
		{Name: "beta", Version: "2.0.0"},
	}

	plan := buildSyncPlan(specs, entries)
	if len(plan.RefreshTargets) != 1 || plan.RefreshTargets[0] != "alpha" {
		t.Fatalf("expected alpha to refresh, got %v", plan.RefreshTargets)
	}
	if plan.SkippedUnpinned != 1 {
		t.Fatalf("expected 1 unpinned skip, got %d", plan.SkippedUnpinned)
	}
	if plan.LocalCount != 1 {
		t.Fatalf("expected 1 local plugin, got %d", plan.LocalCount)
	}

	row := findSyncRow(plan.Rows, "alpha")
	if row.Status != string(model.StatusMismatch) || row.Action != "refresh" {
		t.Fatalf("expected alpha mismatch/refresh, got %s/%s", row.Status, row.Action)
	}

	row = findSyncRow(plan.Rows, "beta")
	if row.Status != string(model.StatusUnknown) || row.Action != "skip" {
		t.Fatalf("expected beta unknown/skip, got %s/%s", row.Status, row.Action)
	}

	row = findSyncRow(plan.Rows, "local-plugin")
	if row.Status != string(model.StatusUnmanaged) || row.Action != "skip" {
		t.Fatalf("expected local unmanaged/skip, got %s/%s", row.Status, row.Action)
	}
}

func TestSyncCommandMissingCacheDir(t *testing.T) {
	root := t.TempDir()
	missing := root + "/missing"
	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    missing,
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := syncCommand(opts, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected error exit code, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "failed to scan cache directory") {
		t.Fatalf("expected cache scan error, got %q", errOut.String())
	}
}

func findSyncRow(rows []syncRow, name string) syncRow {
	for _, row := range rows {
		if row.Name == name {
			return row
		}
	}
	return syncRow{}
}
