package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
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
