package cache

import (
	"context"
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

func Detect(ctx context.Context, cacheDir string) ([]Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	results := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return results, err
		}
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
func Invalidate(ctx context.Context, cacheDir string, pluginName string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if cacheDir == "" {
		return nil, fmt.Errorf("%w: cache directory is empty", ErrInvalidCacheDir)
	}
	info, err := os.Stat(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCacheDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCacheDir, cacheDir)
	}
	if pluginName == "" {
		return nil, fmt.Errorf("plugin name is required")
	}

	entries, err := Detect(ctx, cacheDir)
	if err != nil {
		return nil, fmt.Errorf("detect cache entries: %w", err)
	}

	base, err := filepath.Abs(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("resolve cache directory: %w", err)
	}

	removed := []string{}
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return removed, err
		}
		if entry.Name != pluginName {
			continue
		}
		if err := ensureWithin(base, entry.Path); err != nil {
			return removed, fmt.Errorf("validate cache path %q: %w", entry.Path, err)
		}
		if err := os.RemoveAll(entry.Path); err != nil {
			return removed, fmt.Errorf("remove cache entry %q: %w", entry.Path, err)
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
