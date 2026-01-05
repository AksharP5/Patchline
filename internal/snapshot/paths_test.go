package snapshot

import (
	"path/filepath"
	"testing"
)

func TestResolveDirUsesOverride(t *testing.T) {
	override := t.TempDir()
	resolved, candidates := ResolveDir(override)
	if resolved != override {
		t.Fatalf("expected override %s, got %s", override, resolved)
	}
	if len(candidates) != 0 {
		t.Fatalf("expected no candidates when override is set, got %#v", candidates)
	}
}

func TestCandidateDirsIncludeXDGData(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_DATA_HOME", root)

	candidates := CandidateDirs()
	want := filepath.Join(root, "patchline", "snapshots")
	found := false
	for _, candidate := range candidates {
		if candidate == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected candidate %s, got %#v", want, candidates)
	}
}
