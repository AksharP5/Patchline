package snapshot

import (
	"os"
	"path/filepath"
	"runtime"
)

// ResolveDir returns the snapshot directory and candidate paths.
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
	if len(candidates) > 0 {
		return candidates[0], candidates
	}
	return "", candidates
}

// CandidateDirs returns default snapshot directory candidates.
func CandidateDirs() []string {
	dirs := []string{}
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		dirs = append(dirs, filepath.Join(dataHome, "patchline", "snapshots"))
	}

	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			dirs = append(dirs, filepath.Join(localAppData, "patchline", "snapshots"))
		}
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		switch runtime.GOOS {
		case "darwin":
			dirs = append(dirs, filepath.Join(home, "Library", "Application Support", "patchline", "snapshots"))
		case "windows":
			dirs = append(dirs, filepath.Join(home, "AppData", "Local", "patchline", "snapshots"))
		default:
			dirs = append(dirs, filepath.Join(home, ".local", "share", "patchline", "snapshots"))
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
