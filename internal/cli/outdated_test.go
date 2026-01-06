package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOutdatedCommandOfflineNotesAndLocalPlugins(t *testing.T) {
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("XDG_DATA_HOME", root)
	configPath := filepath.Join(root, "opencode.json")
	config := []byte(`{"plugins":["pkg@1.0.0"]}`)
	if err := os.WriteFile(configPath, config, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	localDir := filepath.Join(root, ".opencode", "plugin")
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatalf("mkdir local plugins: %v", err)
	}
	localPlugin := filepath.Join(localDir, "local.js")
	if err := os.WriteFile(localPlugin, []byte("export default {}"), 0o644); err != nil {
		t.Fatalf("write local plugin: %v", err)
	}

	cacheDir := filepath.Join(root, "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("mkdir cache: %v", err)
	}

	opts := CommonOptions{
		ProjectRoot: root,
		CacheDir:    cacheDir,
		Offline:     true,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := outdatedCommand(opts, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr=%q)", code, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Note: offline mode enabled") {
		t.Fatalf("expected offline note, got %q", output)
	}
	if !strings.Contains(output, "local plugins are unmanaged") {
		t.Fatalf("expected local plugin note, got %q", output)
	}
}
