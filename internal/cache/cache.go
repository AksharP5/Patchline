package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Name    string
	Version string
	Path    string
}

type packageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Detect(cacheDir string) ([]Entry, error) {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	results := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		packagePath := filepath.Join(cacheDir, entry.Name(), "package.json")
		data, err := os.ReadFile(packagePath)
		if err != nil {
			continue
		}

		var pkg packageJSON
		if err := json.Unmarshal(data, &pkg); err != nil {
			continue
		}
		if pkg.Name == "" || pkg.Version == "" {
			continue
		}

		results = append(results, Entry{
			Name:    pkg.Name,
			Version: pkg.Version,
			Path:    filepath.Join(cacheDir, entry.Name()),
		})
	}

	return results, nil
}

// Invalidate removes cached plugin directories that match the npm package name.
func Invalidate(cacheDir string, pluginName string) ([]string, error) {
	if pluginName == "" {
		return nil, fmt.Errorf("plugin name is required")
	}

	entries, err := Detect(cacheDir)
	if err != nil {
		return nil, err
	}

	base, err := filepath.Abs(cacheDir)
	if err != nil {
		return nil, err
	}

	removed := []string{}
	for _, entry := range entries {
		if entry.Name != pluginName {
			continue
		}
		if err := ensureWithin(base, entry.Path); err != nil {
			return removed, err
		}
		if err := os.RemoveAll(entry.Path); err != nil {
			return removed, err
		}
		removed = append(removed, entry.Path)
	}

	return removed, nil
}

func ensureWithin(base string, target string) error {
	base = filepath.Clean(base)
	target = filepath.Clean(target)
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(base, absTarget)
	if err != nil {
		return err
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return ErrUnsafePath
	}
	return nil
}
