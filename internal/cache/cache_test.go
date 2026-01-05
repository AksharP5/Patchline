package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInvalidateRemovesMatchingPlugin(t *testing.T) {
	cacheDir := t.TempDir()
	fooDir := filepath.Join(cacheDir, "foo")
	barDir := filepath.Join(cacheDir, "bar")

	writePackageJSON(t, fooDir, `{"name":"foo","version":"1.0.0"}`)
	writePackageJSON(t, barDir, `{"name":"bar","version":"2.0.0"}`)

	removed, err := Invalidate(cacheDir, "foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removal, got %d", len(removed))
	}
	if _, err := os.Stat(fooDir); !os.IsNotExist(err) {
		t.Fatalf("expected foo dir removed, got %v", err)
	}
	if _, err := os.Stat(barDir); err != nil {
		t.Fatalf("expected bar dir to remain, got %v", err)
	}
}

func TestInvalidateMissingPluginIsNoop(t *testing.T) {
	cacheDir := t.TempDir()
	fooDir := filepath.Join(cacheDir, "foo")
	writePackageJSON(t, fooDir, `{"name":"foo","version":"1.0.0"}`)

	removed, err := Invalidate(cacheDir, "missing")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(removed) != 0 {
		t.Fatalf("expected no removals, got %d", len(removed))
	}
}

func TestEnsureWithin(t *testing.T) {
	base := t.TempDir()
	child := filepath.Join(base, "child")
	if err := ensureWithin(base, child); err != nil {
		t.Fatalf("expected within base, got %v", err)
	}
	if err := ensureWithin(base, base); err == nil {
		t.Fatalf("expected error for base path")
	}
	outside := filepath.Join(base, "..", "other")
	if err := ensureWithin(base, outside); err == nil {
		t.Fatalf("expected error for outside path")
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
