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
