package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

var Version = "dev"

const toolName = "patchline"

type CommonOptions struct {
	ProjectRoot  string
	GlobalConfig string
	CacheDir     string
	SnapshotDir  string
	Offline      bool
	LocalDirs    stringSliceFlag
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	cmd := args[0]
	switch cmd {
	case "list":
		return runList(args[1:], stdout, stderr)
	case "outdated":
		return runOutdated(args[1:], stdout, stderr)
	case "sync":
		return runSync(args[1:], stdout, stderr)
	case "upgrade":
		return runUpgrade(args[1:], stdout, stderr)
	case "rollback":
		return runRollback(args[1:], stdout, stderr)
	case "snapshot":
		return runSnapshot(args[1:], stdout, stderr)
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "%s %s\n", toolName, Version)
		return 0
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", cmd)
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	lines := []string{
		"Usage:",
		"  patchline <command> [options]",
		"",
		"Commands:",
		"  list       Show declared and installed plugins",
		"  outdated   Show plugins with newer versions",
		"  sync       Refresh cache to match pinned config",
		"  upgrade    Pin and refresh plugins to a target version",
		"  rollback   Restore the most recent plugin snapshot",
		"  snapshot   Save a snapshot of current plugin state",
		"  version    Print version information",
	}
	fmt.Fprintln(w, strings.Join(lines, "\n"))
}

func runList(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	opts := bindCommonFlags(fs)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	return listCommand(*opts, stdout, stderr)
}

func runOutdated(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("outdated", flag.ContinueOnError)
	opts := bindCommonFlags(fs)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = opts
	fmt.Fprintln(stderr, "outdated is not implemented yet")
	return 1
}

func runSync(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	opts := bindCommonFlags(fs)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = opts
	fmt.Fprintln(stderr, "sync is not implemented yet")
	return 1
}

func runUpgrade(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	var target string
	var major bool
	var minor bool
	var patch bool
	var all bool
	opts := bindCommonFlags(fs)
	fs.StringVar(&target, "to", "", "explicit target version")
	fs.BoolVar(&major, "major", false, "upgrade to latest major")
	fs.BoolVar(&minor, "minor", false, "upgrade to latest minor")
	fs.BoolVar(&patch, "patch", false, "upgrade to latest patch")
	fs.BoolVar(&all, "all", false, "upgrade all plugins")
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = opts

	name := ""
	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}
	if !all && name == "" {
		fmt.Fprintln(stderr, "missing plugin name or --all")
		return 2
	}
	if all && name != "" {
		fmt.Fprintln(stderr, "cannot use --all with a plugin name")
		return 2
	}
	if flagCount(target != "", major, minor, patch) > 1 {
		fmt.Fprintln(stderr, "only one of --to, --major, --minor, --patch is allowed")
		return 2
	}

	fmt.Fprintln(stderr, "upgrade is not implemented yet")
	return 1
}

func runRollback(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("rollback", flag.ContinueOnError)
	opts := bindCommonFlags(fs)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = opts
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "missing plugin name")
		return 2
	}
	fmt.Fprintln(stderr, "rollback is not implemented yet")
	return 1
}

func runSnapshot(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	opts := bindCommonFlags(fs)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = opts
	fmt.Fprintln(stderr, "snapshot is not implemented yet")
	return 1
}

func bindCommonFlags(fs *flag.FlagSet) *CommonOptions {
	opts := &CommonOptions{}
	fs.StringVar(&opts.ProjectRoot, "project", "", "project root to scan for opencode.json")
	fs.StringVar(&opts.GlobalConfig, "global-config", "", "override global opencode.json path")
	fs.StringVar(&opts.CacheDir, "cache-dir", "", "override OpenCode plugin cache directory")
	fs.StringVar(&opts.SnapshotDir, "snapshot-dir", "", "override snapshot storage directory")
	fs.BoolVar(&opts.Offline, "offline", false, "disable registry network calls")
	fs.Var(&opts.LocalDirs, "local-dir", "additional local plugin directory (repeatable)")
	return opts
}

func flagCount(values ...bool) int {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}
