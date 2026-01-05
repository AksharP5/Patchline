package npm

import (
	"fmt"
	"strconv"
	"strings"
)

// CompareSemver compares two semver strings and returns comparison result.
// The bool return is false when parsing fails.
func CompareSemver(current string, latest string) (int, bool) {
	currentSemver, ok := parseSemver(current)
	if !ok {
		return 0, false
	}
	latestSemver, ok := parseSemver(latest)
	if !ok {
		return 0, false
	}
	return compareSemver(currentSemver, latestSemver), true
}

type Semver struct {
	Major int
	Minor int
	Patch int
}

type UpgradeMode string

const (
	UpgradeLatest UpgradeMode = "latest"
	UpgradeMajor  UpgradeMode = "major"
	UpgradeMinor  UpgradeMode = "minor"
	UpgradePatch  UpgradeMode = "patch"
)

// SelectTargetVersion chooses a version string based on the upgrade mode.
func SelectTargetVersion(latest string, versions []string, base string, mode UpgradeMode) (string, error) {
	switch mode {
	case UpgradeLatest, "":
		if latest == "" {
			return "", fmt.Errorf("latest version is unavailable")
		}
		return latest, nil
	case UpgradeMajor:
		return selectHighest(versions, func(_ Semver) bool { return true })
	case UpgradeMinor:
		baseSemver, ok := parseSemver(base)
		if !ok {
			return "", fmt.Errorf("base version is not semver: %s", base)
		}
		return selectHighest(versions, func(candidate Semver) bool {
			return candidate.Major == baseSemver.Major
		})
	case UpgradePatch:
		baseSemver, ok := parseSemver(base)
		if !ok {
			return "", fmt.Errorf("base version is not semver: %s", base)
		}
		return selectHighest(versions, func(candidate Semver) bool {
			return candidate.Major == baseSemver.Major && candidate.Minor == baseSemver.Minor
		})
	default:
		return "", fmt.Errorf("unknown upgrade mode: %s", mode)
	}
}

func selectHighest(versions []string, accept func(Semver) bool) (string, error) {
	best := ""
	var bestSemver Semver
	found := false
	for _, version := range versions {
		semver, ok := parseSemver(version)
		if !ok || !accept(semver) {
			continue
		}
		if !found || compareSemver(semver, bestSemver) > 0 {
			best = version
			bestSemver = semver
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("no matching versions found")
	}
	return best, nil
}

func parseSemver(version string) (Semver, bool) {
	clean := strings.TrimSpace(version)
	if clean == "" {
		return Semver{}, false
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
		return Semver{}, false
	}

	major, ok := parsePart(parts, 0)
	if !ok {
		return Semver{}, false
	}
	minor, ok := parsePart(parts, 1)
	if !ok {
		return Semver{}, false
	}
	patch, ok := parsePart(parts, 2)
	if !ok {
		return Semver{}, false
	}
	return Semver{Major: major, Minor: minor, Patch: patch}, true
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

func compareSemver(a Semver, b Semver) int {
	if a.Major != b.Major {
		return compareInt(a.Major, b.Major)
	}
	if a.Minor != b.Minor {
		return compareInt(a.Minor, b.Minor)
	}
	if a.Patch != b.Patch {
		return compareInt(a.Patch, b.Patch)
	}
	return 0
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
