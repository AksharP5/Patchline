package opencode

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectConfigWalksParents(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir child: %v", err)
	}

	configPath := filepath.Join(root, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := findProjectConfig(child)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}

func TestFindProjectConfigPrefersNonDotfile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".opencode.json"), []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	configPath := filepath.Join(root, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := findProjectConfig(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}

func TestFindProjectConfigFindsDotfile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := findProjectConfig(root)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}

func TestFindProjectConfigExplicitFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := findProjectConfig(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}

func TestFindProjectConfigExplicitDotfile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := findProjectConfig(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}

func TestLoadPluginSpecsSupportsBothKeys(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	content := []byte(`{"plugin":["alpha@1.0.0"],"plugins":["beta@2.0.0"]}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	plugins, err := loadPluginSpecs(configPath, SourceProject)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(plugins))
	}

	seen := map[string]string{}
	for _, plugin := range plugins {
		seen[plugin.Name] = plugin.Pinned
	}
	if seen["alpha"] != "1.0.0" {
		t.Fatalf("expected alpha@1.0.0, got %q", seen["alpha"])
	}
	if seen["beta"] != "2.0.0" {
		t.Fatalf("expected beta@2.0.0, got %q", seen["beta"])
	}
}

func TestDiscoverLocalPluginsFiltersExtensions(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"local.js":    "export default {}",
		"skip.txt":    "skip",
		"another.cjs": "module.exports = {}",
		"last.mjs":    "export default {}",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	plugins := discoverLocalPlugins([]string{root})
	if len(plugins) != 3 {
		t.Fatalf("expected 3 local plugins, got %d", len(plugins))
	}

	seen := map[string]bool{}
	for _, plugin := range plugins {
		seen[plugin.Name] = true
	}
	if !seen["local"] || !seen["another"] || !seen["last"] {
		t.Fatalf("unexpected local plugin names: %#v", seen)
	}
}

func TestParseSpecScopedPackage(t *testing.T) {
	name, pinned := parseSpec("@scope/pkg@1.2.3")
	if name != "@scope/pkg" || pinned != "1.2.3" {
		t.Fatalf("expected @scope/pkg@1.2.3, got %q@%q", name, pinned)
	}
}

func TestGlobalConfigCandidatesIncludeXDG(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)

	candidates := globalConfigCandidates()
	wantFiles := []string{"opencode.json", ".opencode.json"}
	for _, filename := range wantFiles {
		want := filepath.Join(root, "opencode", filename)
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
}

func TestResolveGlobalConfigPrefersNonDotfile(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)

	configDir := filepath.Join(root, "opencode")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	dotPath := filepath.Join(configDir, ".opencode.json")
	if err := os.WriteFile(dotPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write dot config: %v", err)
	}

	configPath := filepath.Join(configDir, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	found, err := resolveGlobalConfig("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found != configPath {
		t.Fatalf("expected %s, got %s", configPath, found)
	}
}
