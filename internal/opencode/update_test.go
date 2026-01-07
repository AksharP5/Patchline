package opencode

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdatePluginSpecPluginKey(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "opencode.json")
	data := `{
  "plugin": ["alpha@1.0.0", "@scope/beta@2.0.0"],
  "other": "value"
}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := UpdatePluginSpec(path, "alpha", "alpha@1.2.0"); err != nil {
		t.Fatalf("update: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(updated)
	if !strings.Contains(content, "alpha@1.2.0") {
		t.Fatalf("expected updated spec in file")
	}
	if !strings.Contains(content, "@scope/beta@2.0.0") {
		t.Fatalf("expected other plugin preserved")
	}
}

func TestUpdatePluginSpecPluginsKey(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "opencode.json")
	data := `{
  "plugins": ["gamma@3.0.0"]
}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := UpdatePluginSpec(path, "gamma", "gamma@3.1.0"); err != nil {
		t.Fatalf("update: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(updated), "gamma@3.1.0") {
		t.Fatalf("expected updated spec")
	}
}

func TestUpdatePluginSpecJSONC(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "opencode.json")
	data := `// top-level comment
{
  "plugin": ["delta@1.0.0"], // inline comment
  /* block comment */
  "other": "value"
}
`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := UpdatePluginSpec(path, "delta", "delta@1.1.0"); err != nil {
		t.Fatalf("update: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(updated), "delta@1.1.0") {
		t.Fatalf("expected updated spec")
	}
}

func TestUpdatePluginSpecNotFound(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "opencode.json")
	data := `{"plugin": ["alpha@1.0.0"]}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	err := UpdatePluginSpec(path, "missing", "missing@1.0.0")
	if !errors.Is(err, ErrPluginNotFound) {
		t.Fatalf("expected ErrPluginNotFound, got %v", err)
	}
}
