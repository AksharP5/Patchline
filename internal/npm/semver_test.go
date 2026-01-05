package npm

import "testing"

type semverCase struct {
	current  string
	latest   string
	cmp      int
	ok       bool
	caseName string
}

func TestCompareSemver(t *testing.T) {
	cases := []semverCase{
		{caseName: "equal", current: "1.2.3", latest: "1.2.3", cmp: 0, ok: true},
		{caseName: "patch", current: "1.2.3", latest: "1.2.4", cmp: -1, ok: true},
		{caseName: "minor", current: "1.2.3", latest: "1.3.0", cmp: -1, ok: true},
		{caseName: "major", current: "1.2.3", latest: "2.0.0", cmp: -1, ok: true},
		{caseName: "prefix-v", current: "v1.2.3", latest: "1.2.4", cmp: -1, ok: true},
		{caseName: "prerelease", current: "1.2.3-beta.1", latest: "1.2.4", cmp: -1, ok: true},
		{caseName: "build-metadata", current: "1.2.3+build.1", latest: "1.2.3", cmp: 0, ok: true},
		{caseName: "missing-patch", current: "1.2", latest: "1.2.0", cmp: 0, ok: true},
		{caseName: "compare-large", current: "1.10.0", latest: "1.2.0", cmp: 1, ok: true},
		{caseName: "invalid", current: "1.x", latest: "1.2.0", cmp: 0, ok: false},
	}

	for _, tc := range cases {
		cmp, ok := CompareSemver(tc.current, tc.latest)
		if ok != tc.ok {
			t.Fatalf("%s: expected ok=%v, got %v", tc.caseName, tc.ok, ok)
		}
		if ok && cmp != tc.cmp {
			t.Fatalf("%s: expected cmp=%d, got %d", tc.caseName, tc.cmp, cmp)
		}
	}
}

func TestSelectTargetVersion(t *testing.T) {
	versions := []string{"1.2.0", "1.3.5", "2.0.0", "invalid", "1.3.1"}

	latest, err := SelectTargetVersion("2.0.0", versions, "", UpgradeLatest)
	if err != nil || latest != "2.0.0" {
		t.Fatalf("expected latest 2.0.0, got %s (%v)", latest, err)
	}

	major, err := SelectTargetVersion("", versions, "", UpgradeMajor)
	if err != nil || major != "2.0.0" {
		t.Fatalf("expected major 2.0.0, got %s (%v)", major, err)
	}

	minor, err := SelectTargetVersion("", versions, "1.3.0", UpgradeMinor)
	if err != nil || minor != "1.3.5" {
		t.Fatalf("expected minor 1.3.5, got %s (%v)", minor, err)
	}

	patch, err := SelectTargetVersion("", versions, "1.3.0", UpgradePatch)
	if err != nil || patch != "1.3.5" {
		t.Fatalf("expected patch 1.3.5, got %s (%v)", patch, err)
	}
}

func TestSelectTargetVersionMissingBase(t *testing.T) {
	if _, err := SelectTargetVersion("", []string{"1.0.0"}, "", UpgradeMinor); err == nil {
		t.Fatalf("expected error for missing base version")
	}
}
