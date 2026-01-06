package cli

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AksharP5/Patchline/internal/snapshot"
)

func TestSnapshotUpgradeRollbackFlow(t *testing.T) {
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
	if code := snapshotCommand(opts, &stdout, &stderr); code != 0 {
		t.Fatalf("snapshot failed: %d %s", code, stderr.String())
	}

	entries := readSnapshotEntries(t, snapshotDir, "alpha")
	if len(entries) == 0 {
		t.Fatalf("expected snapshot entries")
	}
	if !hasSnapshotSpec(entries, "alpha@1.0.0") {
		t.Fatalf("expected snapshot for alpha@1.0.0")
	}

	stdout.Reset()
	stderr.Reset()
	if code := upgradeCommand(opts, "alpha", "1.1.0", "", false, &stdout, &stderr); code != 0 {
		t.Fatalf("upgrade failed: %d %s", code, stderr.String())
	}

	updated, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(updated), "alpha@1.1.0") {
		t.Fatalf("expected upgrade config, got %s", string(updated))
	}
	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Fatalf("expected cache invalidation, got %v", err)
	}

	writePackageJSON(t, pluginDir, `{"name":"alpha","version":"1.1.0"}`)

	stdout.Reset()
	stderr.Reset()
	if code := rollbackCommand(opts, "alpha", &stdout, &stderr); code != 0 {
		t.Fatalf("rollback failed: %d %s", code, stderr.String())
	}

	restored, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(restored), "alpha@1.0.0") {
		t.Fatalf("expected rollback config, got %s", string(restored))
	}
	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Fatalf("expected cache invalidation after rollback, got %v", err)
	}

	entries = readSnapshotEntries(t, snapshotDir, "alpha")
	if len(entries) < 2 {
		t.Fatalf("expected multiple snapshot entries, got %d", len(entries))
	}
}

func readSnapshotEntries(t *testing.T, dir string, pluginName string) []snapshot.Entry {
	t.Helper()
	path := filepath.Join(dir, url.PathEscape(pluginName)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	var entries []snapshot.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	return entries
}

func hasSnapshotSpec(entries []snapshot.Entry, spec string) bool {
	for _, entry := range entries {
		if entry.PreviousSpec == spec {
			return true
		}
	}
	return false
}
