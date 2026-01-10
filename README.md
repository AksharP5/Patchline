# Patchline

Streamline OpenCode plugin updates by inspecting config and cache state, pinning versions, and invalidating stale cache entries.

Patchline does not install plugins or run OpenCode. After running Patchline, launch OpenCode to reinstall or refresh plugins.

## Install

### Homebrew

```
brew install AksharP5/tap/patchline
```

### From source

```
go install github.com/AksharP5/Patchline/cmd/patchline@latest
```

### npm

```
npm i -g patchline
```

## Usage

```
patchline list
patchline outdated
patchline sync
patchline upgrade <plugin> --to 1.2.3
patchline upgrade <plugin> --major|--minor|--patch
patchline upgrade --all --minor
patchline snapshot <plugin>
patchline rollback <plugin>
patchline version
```

## Common flags

- `--project <dir>`: project root to scan for `opencode.json`.
- `--global-config <file>`: override the global config path.
- `--cache-dir <dir>`: override the OpenCode plugin cache directory.
- `--snapshot-dir <dir>`: override where snapshots are stored.
- `--local-dir <dir>`: add an extra local plugin directory (repeatable).
- `--offline`: skip npm registry calls.

## Status meanings

- `missing`: declared in config, but no cache entry was found.
- `mismatch`: cache entry exists but does not match the pinned version.
- `outdated`: installed version is behind the npm registry latest.
- `local/unmanaged`: plugin is a local file and not managed by npm.

## Troubleshooting

- If plugins show as `missing`, run OpenCode to install them.
- If Patchline cannot find your config or cache, pass `--global-config` or `--cache-dir`.
