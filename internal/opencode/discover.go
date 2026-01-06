package opencode

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var localPluginExtensions = map[string]bool{
	".js":  true,
	".cjs": true,
	".mjs": true,
}

var configFilenames = []string{"opencode.json", ".opencode.json"}

func Discover(projectRoot string, globalConfigPath string, localDirs []string) (DiscoveryResult, error) {
	result := DiscoveryResult{}

	projectPath, err := findProjectConfig(projectRoot)
	if err != nil && !errors.Is(err, ErrConfigNotFound) {
		return result, err
	}
	if projectPath != "" {
		plugins, err := loadPluginSpecs(projectPath, SourceProject)
		if err != nil {
			return result, err
		}
		result.Plugins = append(result.Plugins, plugins...)
	}

	globalPath, err := resolveGlobalConfig(globalConfigPath)
	if err != nil && !errors.Is(err, ErrConfigNotFound) {
		return result, err
	}
	if globalPath != "" {
		plugins, err := loadPluginSpecs(globalPath, SourceGlobal)
		if err != nil {
			return result, err
		}
		result.Plugins = append(result.Plugins, plugins...)
	}

	localCandidates := append([]string{}, localDirs...)
	localCandidates = append(localCandidates, defaultLocalPluginDirs(projectRoot, globalPath)...)
	result.Plugins = append(result.Plugins, discoverLocalPlugins(localCandidates)...)

	return result, nil
}

func loadPluginSpecs(path string, source Source) ([]PluginSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	declared := append([]string{}, cfg.Plugins...)
	declared = append(declared, cfg.PluginsAlt...)

	plugins := make([]PluginSpec, 0, len(declared))
	for _, spec := range declared {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		name, pinned := parseSpec(spec)
		if name == "" {
			continue
		}
		plugins = append(plugins, PluginSpec{
			Name:         name,
			DeclaredSpec: spec,
			Pinned:       pinned,
			Source:       source,
			ConfigPath:   path,
		})
	}

	return plugins, nil
}

func parseSpec(spec string) (string, string) {
	at := strings.LastIndex(spec, "@")
	if at <= 0 {
		return spec, ""
	}

	name := spec[:at]
	version := spec[at+1:]
	if name == "" || version == "" {
		return spec, ""
	}
	return name, version
}

func findProjectConfig(projectRoot string) (string, error) {
	root := projectRoot
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		root = cwd
	}

	info, err := os.Stat(root)
	if err == nil && !info.IsDir() {
		if isConfigFilename(filepath.Base(root)) {
			return root, nil
		}
		return "", fmt.Errorf("project path is not a directory: %s", root)
	}

	root = filepath.Clean(root)
	for {
		for _, name := range configFilenames {
			candidate := filepath.Join(root, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}
		parent := filepath.Dir(root)
		if parent == root {
			return "", ErrConfigNotFound
		}
		root = parent
	}
}

func resolveGlobalConfig(override string) (string, error) {
	if override != "" {
		if fileExists(override) {
			return override, nil
		}
		return "", fmt.Errorf("%w: %s", ErrConfigNotFound, override)
	}

	for _, candidate := range globalConfigCandidates() {
		if fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", ErrConfigNotFound
}

func globalConfigCandidates() []string {
	paths := []string{}
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		for _, name := range configFilenames {
			paths = append(paths, filepath.Join(configHome, "opencode", name))
		}
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		var baseDirs []string
		switch runtime.GOOS {
		case "darwin":
			baseDirs = append(baseDirs, filepath.Join(home, "Library", "Application Support", "opencode"))
		default:
			baseDirs = append(baseDirs, filepath.Join(home, ".config", "opencode"))
		}
		baseDirs = append(baseDirs, filepath.Join(home, ".opencode"))
		for _, base := range baseDirs {
			for _, name := range configFilenames {
				paths = append(paths, filepath.Join(base, name))
			}
		}
	}

	return uniqueStrings(paths)
}

func defaultLocalPluginDirs(projectRoot string, globalConfigPath string) []string {
	paths := []string{}
	if globalConfigPath != "" {
		paths = append(paths, filepath.Join(filepath.Dir(globalConfigPath), "plugins"))
	}

	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		paths = append(paths, filepath.Join(dataHome, "opencode", "plugins"))
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		switch runtime.GOOS {
		case "darwin":
			paths = append(paths, filepath.Join(home, "Library", "Application Support", "opencode", "plugins"))
		default:
			paths = append(paths, filepath.Join(home, ".local", "share", "opencode", "plugins"))
		}
	}

	if projectRoot != "" {
		paths = append(paths,
			filepath.Join(projectRoot, ".opencode", "plugins"),
			filepath.Join(projectRoot, "opencode", "plugins"),
		)
	}

	return uniqueStrings(paths)
}

func discoverLocalPlugins(dirs []string) []PluginSpec {
	plugins := []PluginSpec{}
	for _, dir := range uniqueStrings(dirs) {
		if !dirExists(dir) {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if !localPluginExtensions[ext] {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), ext)
			path := filepath.Join(dir, entry.Name())
			plugins = append(plugins, PluginSpec{
				Name:         name,
				DeclaredSpec: path,
				Source:       SourceLocal,
				LocalPath:    path,
			})
		}
	}

	return plugins
}

func isConfigFilename(name string) bool {
	for _, candidate := range configFilenames {
		if name == candidate {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
