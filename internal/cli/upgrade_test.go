package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AksharP5/Patchline/internal/opencode"
	"github.com/AksharP5/Patchline/internal/snapshot"
)

func TestSelectUpgradeTargetsPreferProject(t *testing.T) {
	specs := []opencode.PluginSpec{
		{Name: "alpha", ConfigPath: "/project/opencode.json", Source: opencode.SourceProject},
		{Name: "alpha", ConfigPath: "/global/opencode.json", Source: opencode.SourceGlobal},
	}

	targets := selectUpgradeTargets(specs, "alpha", false)
	if len(targets) != 1 {
		t.Fatalf("expected one target, got %d", len(targets))
	}
	if targets[0].ConfigPath != "/project/opencode.json" {
		t.Fatalf("expected project config, got %s", targets[0].ConfigPath)
	}
}

func TestSelectUpgradeTargetsAll(t *testing.T) {
	specs := []opencode.PluginSpec{
		{Name: "alpha", ConfigPath: "/project/opencode.json", Source: opencode.SourceProject},
		{Name: "beta", ConfigPath: "/global/opencode.json", Source: opencode.SourceGlobal},
	}

	targets := selectUpgradeTargets(specs, "", true)
	if len(targets) != 2 {
		t.Fatalf("expected two targets, got %d", len(targets))
	}
}

func TestUpgradeCommandExplicitTarget(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	config := `{"plugin": ["alpha@1.0.0"]}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cacheDir := filepath.Join(root, "cache")
	pluginDir := filepath.Join(cacheDir, "alpha-cache")
	writePackageJSON(t, pluginDir, `{"name":"alpha","version":"1.0.0"}`)

	snapshotDir := filepath.Join(root, "snapshots")
	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    cacheDir,
		SnapshotDir: snapshotDir,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := upgradeCommand(opts, "alpha", "1.2.0", "", false, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}

	updated, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(updated), "alpha@1.2.0") {
		t.Fatalf("expected updated spec, got %s", string(updated))
	}

	store := snapshot.Store{Directory: snapshotDir}
	entry, err := store.Latest("alpha")
	if err != nil {
		t.Fatalf("latest snapshot: %v", err)
	}
	if entry.PreviousSpec != "alpha@1.0.0" {
		t.Fatalf("expected previous spec, got %s", entry.PreviousSpec)
	}

	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Fatalf("expected cache removal, got %v", err)
	}
}

func TestUpgradeCommandSkipsIfUpToDate(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	config := `{"plugin": ["alpha@1.2.0"]}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cacheDir := filepath.Join(root, "cache")
	pluginDir := filepath.Join(cacheDir, "alpha-cache")
	writePackageJSON(t, pluginDir, `{"name":"alpha","version":"1.2.0"}`)

	snapshotDir := filepath.Join(root, "snapshots")
	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    cacheDir,
		SnapshotDir: snapshotDir,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := upgradeCommand(opts, "alpha", "1.2.0", "", false, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Already on alpha@1.2.0") {
		t.Fatalf("expected skip message, got %s", stdout.String())
	}

	store := snapshot.Store{Directory: snapshotDir}
	if _, err := store.Latest("alpha"); !errors.Is(err, snapshot.ErrSnapshotNotFound) {
		t.Fatalf("expected no snapshot, got %v", err)
	}

	if _, err := os.Stat(pluginDir); err != nil {
		t.Fatalf("expected cache retained, got %v", err)
	}
}

func TestChooseBaseVersion(t *testing.T) {
	cases := []struct {
		pinned    string
		installed string
		expected  string
	}{
		{pinned: "1.2.3", installed: "1.0.0", expected: "1.2.3"},
		{pinned: "latest", installed: "1.0.0", expected: "1.0.0"},
		{pinned: "", installed: "1.0.0", expected: "1.0.0"},
		{pinned: "latest", installed: "", expected: ""},
	}

	for _, tc := range cases {
		if got := chooseBaseVersion(tc.pinned, tc.installed); got != tc.expected {
			t.Fatalf("pinned=%s installed=%s: expected %s, got %s", tc.pinned, tc.installed, tc.expected, got)
		}
	}
}
