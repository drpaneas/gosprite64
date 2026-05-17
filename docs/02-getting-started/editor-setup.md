# Editor Setup

## Cursor / VS Code

If `gopls` reports `embedded/*` packages as missing or does not recognize files guarded by `//go:build n64`, configure the workspace so the editor uses the same build tag and toolchain environment as the terminal.

Create `.vscode/settings.json` in the repository with:

```json
{
  "go.buildTags": "n64",
  "go.toolsEnvVars": {
    "GOENV": "${workspaceFolder}/n64.env"
  },
  "gopls": {
    "build.buildFlags": ["-tags=n64"],
    "env": {
      "GOENV": "${workspaceFolder}/n64.env"
    }
  }
}
```

If the Go extension still invokes the wrong `go` binary for editor actions, add:

```json
"go.alternateTools": {
  "go": "go1.24.5-embedded"
}
```

After changing the settings, run `Go: Restart Language Server` or `Developer: Reload Window`.
