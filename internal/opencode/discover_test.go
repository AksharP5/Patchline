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

func TestLoadPluginSpecsJSONC(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	content := []byte(`// comment
{
  "plugin": ["alpha@1.0.0",], /* inline */
  "plugins": ["beta@2.0.0"],
}
`)
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
}

func TestDiscoverLocalPluginsFiltersExtensions(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"local.js":    "export default {}",
		"skip.txt":    "skip",
		"another.cjs": "module.exports = {}",
		"last.mjs":    "export default {}",
		"typed.ts":    "export default {}",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	plugins := discoverLocalPlugins([]string{root})
	if len(plugins) != 4 {
		t.Fatalf("expected 4 local plugins, got %d", len(plugins))
	}

	seen := map[string]bool{}
	for _, plugin := range plugins {
		seen[plugin.Name] = true
	}
	if !seen["local"] || !seen["another"] || !seen["last"] || !seen["typed"] {
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
	want := filepath.Join(root, "opencode", "opencode.json")
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

func TestDefaultLocalPluginDirs(t *testing.T) {
	root := t.TempDir()
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	dirs := defaultLocalPluginDirs(root, "")
	wantGlobal := filepath.Join(configHome, "opencode", "plugin")
	wantProject := filepath.Join(root, ".opencode", "plugin")

	foundGlobal := false
	foundProject := false
	for _, dir := range dirs {
		if dir == wantGlobal {
			foundGlobal = true
		}
		if dir == wantProject {
			foundProject = true
		}
	}
	if !foundGlobal || !foundProject {
		t.Fatalf("expected plugin dirs, got %#v", dirs)
	}
}

func TestDiscoverUsesCustomConfigFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	customPath := filepath.Join(root, "custom.json")
	if err := os.WriteFile(customPath, []byte(`{"plugins":["foo@2.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write custom config: %v", err)
	}
	t.Setenv("OPENCODE_CONFIG", customPath)

	result, err := Discover(root, "", nil)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}

	foundCustom := false
	for _, plugin := range result.Plugins {
		if plugin.Name == "foo" && plugin.Source == SourceCustom && plugin.Pinned == "2.0.0" {
			foundCustom = true
			break
		}
	}
	if !foundCustom {
		t.Fatalf("expected custom config plugin override")
	}
}

func TestDiscoverUsesCustomConfigDir(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "opencode.json")
	if err := os.WriteFile(configPath, []byte(`{"plugins":["foo@1.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	customDir := filepath.Join(root, "custom")
	if err := os.MkdirAll(filepath.Join(customDir, "plugin"), 0o755); err != nil {
		t.Fatalf("mkdir custom plugin dir: %v", err)
	}
	customConfig := filepath.Join(customDir, "opencode.json")
	if err := os.WriteFile(customConfig, []byte(`{"plugins":["bar@2.0.0"]}`), 0o644); err != nil {
		t.Fatalf("write custom config: %v", err)
	}
	localPlugin := filepath.Join(customDir, "plugin", "local.js")
	if err := os.WriteFile(localPlugin, []byte("export default {}"), 0o644); err != nil {
		t.Fatalf("write local plugin: %v", err)
	}
	t.Setenv("OPENCODE_CONFIG_DIR", customDir)

	result, err := Discover(root, "", nil)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}

	foundCustom := false
	foundLocal := false
	for _, plugin := range result.Plugins {
		if plugin.Name == "bar" && plugin.Source == SourceCustomDir {
			foundCustom = true
		}
		if plugin.Name == "local" && plugin.Source == SourceLocal {
			foundLocal = true
		}
	}
	if !foundCustom || !foundLocal {
		t.Fatalf("expected custom config dir and local plugin")
	}
}

func TestGlobalConfigCandidatesIncludeWindowsEnv(t *testing.T) {
	if os.Getenv("OS") != "Windows_NT" {
		t.Skip("windows only")
	}

	root := t.TempDir()
	appData := filepath.Join(root, "roaming")
	localAppData := filepath.Join(root, "local")
	t.Setenv("APPDATA", appData)
	t.Setenv("LOCALAPPDATA", localAppData)
	t.Setenv("XDG_CONFIG_HOME", "")

	candidates := globalConfigCandidates()
	wantApp := filepath.Join(appData, "opencode", "opencode.json")
	wantLocal := filepath.Join(localAppData, "opencode", "opencode.json")

	foundApp := false
	foundLocal := false
	for _, candidate := range candidates {
		if candidate == wantApp {
			foundApp = true
		}
		if candidate == wantLocal {
			foundLocal = true
		}
	}
	if !foundApp {
		t.Fatalf("expected candidate %s, got %#v", wantApp, candidates)
	}
	if !foundLocal {
		t.Fatalf("expected candidate %s, got %#v", wantLocal, candidates)
	}
}

func TestDefaultLocalPluginDirsIncludeWindowsEnv(t *testing.T) {
	if os.Getenv("OS") != "Windows_NT" {
		t.Skip("windows only")
	}

	root := t.TempDir()
	appData := filepath.Join(root, "roaming")
	localAppData := filepath.Join(root, "local")
	t.Setenv("APPDATA", appData)
	t.Setenv("LOCALAPPDATA", localAppData)
	t.Setenv("XDG_CONFIG_HOME", "")

	dirs := defaultLocalPluginDirs(root, "")
	wantApp := filepath.Join(appData, "opencode", "plugin")
	wantLocal := filepath.Join(localAppData, "opencode", "plugin")

	foundApp := false
	foundLocal := false
	for _, dir := range dirs {
		if dir == wantApp {
			foundApp = true
		}
		if dir == wantLocal {
			foundLocal = true
		}
	}
	if !foundApp {
		t.Fatalf("expected plugin dir %s, got %#v", wantApp, dirs)
	}
	if !foundLocal {
		t.Fatalf("expected plugin dir %s, got %#v", wantLocal, dirs)
	}
}
