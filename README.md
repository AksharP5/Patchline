# Patchline

Streamline OpenCode plugin updates.

## Install

### Homebrew

```
brew install AksharP5/tap/patchline
```

### From source

```
go install github.com/AksharP5/Patchline/cmd/patchline@latest
```

## Release

- Push a tag like `v0.1.0` to trigger the GoReleaser workflow.
- Ensure the `HOMEBREW_TAP_GITHUB_TOKEN` secret exists to update `AksharP5/tap`.
