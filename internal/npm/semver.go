package npm

import (
	"strconv"
	"strings"
)

func CompareSemver(current string, latest string) (int, bool) {
	curMajor, curMinor, curPatch, ok := parseSemver(current)
	if !ok {
		return 0, false
	}
	latMajor, latMinor, latPatch, ok := parseSemver(latest)
	if !ok {
		return 0, false
	}

	if curMajor != latMajor {
		return compareInt(curMajor, latMajor), true
	}
	if curMinor != latMinor {
		return compareInt(curMinor, latMinor), true
	}
	if curPatch != latPatch {
		return compareInt(curPatch, latPatch), true
	}
	return 0, true
}

func parseSemver(version string) (int, int, int, bool) {
	clean := strings.TrimSpace(version)
	if clean == "" {
		return 0, 0, 0, false
	}
	clean = strings.TrimPrefix(clean, "v")
	if idx := strings.Index(clean, "+"); idx >= 0 {
		clean = clean[:idx]
	}
	if idx := strings.Index(clean, "-"); idx >= 0 {
		clean = clean[:idx]
	}

	parts := strings.Split(clean, ".")
	if len(parts) == 0 {
		return 0, 0, 0, false
	}

	major, ok := parsePart(parts, 0)
	if !ok {
		return 0, 0, 0, false
	}
	minor, ok := parsePart(parts, 1)
	if !ok {
		return 0, 0, 0, false
	}
	patch, ok := parsePart(parts, 2)
	if !ok {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}

func parsePart(parts []string, index int) (int, bool) {
	if index >= len(parts) {
		return 0, true
	}
	value := strings.TrimSpace(parts[index])
	if value == "" {
		return 0, false
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return number, true
}

func compareInt(a int, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
