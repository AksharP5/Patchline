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
		dirs = append(dirs, filepath.Join(cacheHome, "opencode", "node_modules"))
	}

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			dirs = append(dirs, filepath.Join(localAppData, "opencode", "node_modules"))
		}
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		switch runtime.GOOS {
		case "windows":
			dirs = append(dirs, filepath.Join(home, "AppData", "Local", "opencode", "node_modules"))
		default:
			dirs = append(dirs, filepath.Join(home, ".cache", "opencode", "node_modules"))
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
