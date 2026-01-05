package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AksharP5/Patchline/internal/snapshot"
)

func TestSnapshotCommandWritesEntries(t *testing.T) {
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
	code := snapshotCommand(opts, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}

	store := snapshot.Store{Directory: snapshotDir}
	entry, err := store.Latest("alpha")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if entry.PreviousSpec != "alpha@1.0.0" {
		t.Fatalf("expected previous spec, got %s", entry.PreviousSpec)
	}
	if entry.PreviousInstalled != "1.0.0" {
		t.Fatalf("expected installed version, got %s", entry.PreviousInstalled)
	}
	if entry.ConfigPath != configPath {
		t.Fatalf("expected config path, got %s", entry.ConfigPath)
	}
}

func TestRollbackCommandRestoresConfigAndCache(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	config := `{"plugin": ["alpha@2.0.0"]}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cacheDir := filepath.Join(root, "cache")
	pluginDir := filepath.Join(cacheDir, "alpha-cache")
	writePackageJSON(t, pluginDir, `{"name":"alpha","version":"2.0.0"}`)

	snapshotDir := filepath.Join(root, "snapshots")
	store := snapshot.Store{Directory: snapshotDir}
	err := store.Save(snapshot.Entry{
		PluginName:        "alpha",
		PreviousSpec:      "alpha@1.0.0",
		PreviousInstalled: "1.0.0",
		ConfigPath:        configPath,
		Source:            "project",
		Reason:            "upgrade",
		Timestamp:         time.Now(),
	})
	if err != nil {
		t.Fatalf("save snapshot: %v", err)
	}

	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    cacheDir,
		SnapshotDir: snapshotDir,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := rollbackCommand(opts, "alpha", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}

	updated, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(updated), "alpha@1.0.0") {
		t.Fatalf("expected config rollback, got %s", string(updated))
	}
	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Fatalf("expected cache removal, got %v", err)
	}
}

func writePackageJSON(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}
