package cache

import (
	"context"
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

	removed, err := Invalidate(context.Background(), cacheDir, "foo")
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

	removed, err := Invalidate(context.Background(), cacheDir, "missing")
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

func TestInvalidateEmptyCacheDir(t *testing.T) {
	if _, err := Invalidate(context.Background(), "", "foo"); err == nil {
		t.Fatalf("expected error for empty cache dir")
	}
}

func TestDetectSupportsScopedPackages(t *testing.T) {
	cacheDir := t.TempDir()
	writePackageJSON(t, filepath.Join(cacheDir, "alpha"), `{"name":"alpha","version":"1.0.0"}`)
	writePackageJSON(t, filepath.Join(cacheDir, "@scope", "beta"), `{"name":"@scope/beta","version":"2.0.0"}`)
	if err := os.MkdirAll(filepath.Join(cacheDir, ".bin"), 0o755); err != nil {
		t.Fatalf("mkdir .bin: %v", err)
	}

	entries, err := Detect(context.Background(), cacheDir)
	if err != nil {
		t.Fatalf("detect: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	seen := map[string]bool{}
	for _, entry := range entries {
		seen[entry.Name] = true
	}
	if !seen["alpha"] || !seen["@scope/beta"] {
		t.Fatalf("unexpected entries: %#v", seen)
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
