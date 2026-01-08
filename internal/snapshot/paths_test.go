package snapshot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setHomeEnv(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if volume := filepath.VolumeName(home); volume != "" {
		t.Setenv("HOMEDRIVE", volume)
		t.Setenv("HOMEPATH", strings.TrimPrefix(home, volume))
	}
}

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

func TestResolveDirFindsExistingCandidate(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_DATA_HOME", root)
	setHomeEnv(t, filepath.Join(root, "home"))

	snapshotDir := filepath.Join(root, "patchline", "snapshots")
	if err := os.MkdirAll(snapshotDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	resolved, candidates := ResolveDir("")
	if resolved != snapshotDir {
		t.Fatalf("expected resolved dir %s, got %s", snapshotDir, resolved)
	}
	if len(candidates) == 0 {
		t.Fatalf("expected candidates, got %#v", candidates)
	}
}

func TestResolveDirReturnsFirstCandidateWhenMissing(t *testing.T) {
	root := t.TempDir()
	xdgHome := filepath.Join(root, "xdg")
	home := filepath.Join(root, "home")
	t.Setenv("XDG_DATA_HOME", xdgHome)
	setHomeEnv(t, home)

	resolved, candidates := ResolveDir("")
	want := filepath.Join(xdgHome, "patchline", "snapshots")
	if resolved != want {
		t.Fatalf("expected resolved dir %s, got %s", want, resolved)
	}
	if len(candidates) == 0 {
		t.Fatalf("expected candidates, got %#v", candidates)
	}
}
