package cache

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

func TestCandidateDirsIncludeXDGCache(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", root)

	candidates := CandidateDirs()
	want := filepath.Join(root, "opencode", "node_modules")
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
	t.Setenv("XDG_CACHE_HOME", root)
	setHomeEnv(t, filepath.Join(root, "home"))

	cacheDir := filepath.Join(root, "opencode", "node_modules")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	resolved, candidates := ResolveDir("")
	if resolved != cacheDir {
		t.Fatalf("expected resolved dir %s, got %s", cacheDir, resolved)
	}
	if len(candidates) == 0 {
		t.Fatalf("expected candidates, got %#v", candidates)
	}
}

func TestResolveDirReturnsEmptyWhenNoCandidatesExist(t *testing.T) {
	root := t.TempDir()
	xdgHome := filepath.Join(root, "xdg")
	home := filepath.Join(root, "home")
	t.Setenv("XDG_CACHE_HOME", xdgHome)
	setHomeEnv(t, home)

	resolved, candidates := ResolveDir("")
	if resolved != "" {
		t.Fatalf("expected empty resolved dir, got %s", resolved)
	}
	if len(candidates) == 0 {
		t.Fatalf("expected candidates, got %#v", candidates)
	}

	want := filepath.Join(xdgHome, "opencode", "node_modules")
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
