package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestStringSliceFlagSetAndString(t *testing.T) {
	var values stringSliceFlag
	if err := values.Set("alpha"); err != nil {
		t.Fatalf("set alpha: %v", err)
	}
	if err := values.Set("beta"); err != nil {
		t.Fatalf("set beta: %v", err)
	}
	if got := values.String(); got != "alpha,beta" {
		t.Fatalf("expected joined values, got %q", got)
	}
}

func TestRunVersion(t *testing.T) {
	prev := Version
	Version = "test"
	t.Cleanup(func() { Version = prev })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d", code)
	}
	if !strings.Contains(stdout.String(), "patchline test") {
		t.Fatalf("expected version output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected usage output, got %q", stdout.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"nope"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected error exit code, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected unknown command message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("expected usage output, got %q", stderr.String())
	}
}

func TestRunNoArgs(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected error exit code, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("expected usage output, got %q", stderr.String())
	}
}

func TestRunUpgradeValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "missing name",
			args: []string{},
			want: "missing plugin name",
		},
		{
			name: "all with name",
			args: []string{"--all", "alpha"},
			want: "cannot use --all",
		},
		{
			name: "conflicting flags",
			args: []string{"--major", "--minor", "alpha"},
			want: "only one of",
		},
		{
			name: "all with target",
			args: []string{"--all", "--to", "1.2.3"},
			want: "cannot use --to",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := runUpgrade(tc.args, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("expected error exit code, got %d", code)
			}
			if !strings.Contains(stderr.String(), tc.want) {
				t.Fatalf("expected %q in stderr, got %q", tc.want, stderr.String())
			}
		})
	}
}

func TestRunCommandFlagErrors(t *testing.T) {
	cases := []struct {
		name string
		run  func([]string, io.Writer, io.Writer) int
	}{
		{
			name: "list",
			run:  runList,
		},
		{
			name: "outdated",
			run:  runOutdated,
		},
		{
			name: "sync",
			run:  runSync,
		},
		{
			name: "snapshot",
			run:  runSnapshot,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := tc.run([]string{"--unknown"}, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("expected error exit code, got %d", code)
			}
			if stderr.Len() == 0 {
				t.Fatalf("expected stderr output for flag error")
			}
		})
	}
}

func TestRunRollbackMissingName(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runRollback([]string{}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected error exit code, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing plugin name") {
		t.Fatalf("expected missing plugin name, got %q", stderr.String())
	}
}
