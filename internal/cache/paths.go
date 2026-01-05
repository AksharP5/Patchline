package cache

import (
	"os"
	"path/filepath"
	"runtime"
)

func ResolveDir(override string) (string, []string) {
	if override != "" {
		return override, nil
	}

	candidates := CandidateDirs()
	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return dir, candidates
		}
	}
	return "", candidates
}

func CandidateDirs() []string {
	dirs := []string{}
	if cacheHome := os.Getenv("XDG_CACHE_HOME"); cacheHome != "" {
		dirs = append(dirs,
			filepath.Join(cacheHome, "opencode", "plugins"),
			filepath.Join(cacheHome, "opencode"),
		)
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		switch runtime.GOOS {
		case "darwin":
			dirs = append(dirs,
				filepath.Join(home, "Library", "Caches", "opencode", "plugins"),
				filepath.Join(home, "Library", "Caches", "opencode"),
			)
		default:
			dirs = append(dirs,
				filepath.Join(home, ".cache", "opencode", "plugins"),
				filepath.Join(home, ".cache", "opencode"),
			)
		}
	}

	return uniqueStrings(dirs)
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
