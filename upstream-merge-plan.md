# Upstream Merge Plan

## Context

We maintain a fork of `tomasz-tomczyk/crit` with two custom changes that need to survive each upstream merge:

1. **Files outside repo root fix** — When crit is given a file outside the repo root (e.g. `~/.claude/plans/foo.md`), `filepath.Rel` produces a `../..` path that the `/files/` endpoint rejects. Our fix: use absolute paths and validate against known session files via `IsSessionFile()`. Upstream does NOT have this fix.
   - Method lives in `internal/session/session.go` as `IsSessionFile` (exported, uppercase — required for cross-package access from `internal/server/`)
   - Validation block lives in `internal/server/server.go` `handleFiles` function
2. **No Homebrew tap push** — upstream's release.yml pushes to `tomasz-tomczyk/homebrew-tap` using a `HOMEBREW_TAP_TOKEN` secret we don't have. Strip this step after merge if present.

### Removed custom changes (no longer needed)
- **`--global` flag for `crit install`** — This flag was removed upstream. Global install is now determined automatically: if `cwd == $HOME`, the install is global. Do NOT re-add `--global`.

### Upstream structural notes (as of v0.16.3+)
- **`integrations/` moved to repo root** — was `cmd/crit/integrations/`, now `integrations/` at the repo root, with `integrations/embed.go` providing the Go embed.
- **claude-code integration is now SKILLS** — `integrations/claude-code/skills/crit/SKILL.md` and `integrations/claude-code/skills/crit-cli/SKILL.md`, not a single `.claude/commands/crit.md` file.
- **`integrations.go` renamed to `cli_install.go`** — the install CLI subcommand.
- **Package layout refactored** — code moved from flat `cmd/crit/` into `internal/` packages (server, session, vcs, github, share, etc.) and `web/` for frontend assets. `cmd/crit/` is now a thin CLI wiring layer.

Items that are NOT concerns:
- The Nix flake and marketplace JSON files reference `tomasz-tomczyk/crit` but are inert in our fork.
- `install.sh` — our cross-platform install script, local-only, just keep it.

## Merge History

| Version | Date | Conflicts | Notes |
|---------|------|-----------|-------|
| v0.8.3 | 2026-04-10 | main.go, session.go, server.go | Initial merge from v0.7.0 |
| v0.9.2 | 2026-04-17 | main.go (help text), session.go (isSessionFile) | Clean merge, 20 upstream commits |
| v0.16.2 | 2026-05-xx | Various (large refactor range) | Merged through v0.16.2 |
| v0.16.5 | 2026-06-29 | internal/config/config.go (OpenCmd added; kept ShareURL=""), test/e2e/tests/settings-panel.spec.ts (config card assertions) | Massive repo layout refactor (v0.16.3): internal/ packages, web/, integrations/ moved to root; 127 Go files needed module path updated from tomasz-tomczyk to JoshEllinger; isSessionFile exported as IsSessionFile for cross-package use |

## Steps for future merges

### 1. Merge upstream/main (or tag)

```bash
git fetch upstream --tags
git merge v0.X.Y
```

### 2. Resolve conflicts

Expected conflict areas (our custom code vs upstream changes):

- **internal/config/config.go** — `defaultConfig()` has `ShareURL: ""` (fork intentionally suppresses the crit.md default). Upstream may add new fields (e.g. `OpenCmd`). Take new upstream fields, keep `ShareURL: ""`.
- **internal/session/session.go** — Our `IsSessionFile()` method. Keep it; upstream doesn't have it.
- **internal/server/server.go** — Our `filepath.IsAbs` block in `handleFiles`. Keep it; upstream doesn't have it.

### 3. Verify custom changes survived

```bash
# Files-outside-root fix
grep -rn 'IsSessionFile' internal/session/ internal/server/
grep -rn 'filepath.IsAbs(reqPath)' internal/server/

# No homebrew tap push
grep -n 'homebrew\|HOMEBREW_TAP' .github/workflows/release.yml  # expect NO output
```

### 4. Module import path fix (if new files were added by upstream)

Upstream uses `github.com/tomasz-tomczyk/crit/...`; our fork uses `github.com/JoshEllinger/crit/...`. After any merge that brings in new Go files:

```bash
# Check for stale upstream import paths
grep -r 'github.com/tomasz-tomczyk/crit' . --include='*.go' -l

# Fix them all at once if any are found
find . -name '*.go' -exec grep -l 'github.com/tomasz-tomczyk/crit' {} \; | \
  xargs sed -i '' 's|github.com/tomasz-tomczyk/crit|github.com/JoshEllinger/crit|g'
```

### 5. Build, test, format

```bash
go build -o /tmp/crit-build ./cmd/crit; echo "build $?"
go test ./... 2>&1 | tail -30; echo "test $?"
gofmt -l .         # should be empty
gofmt -w $(gofmt -l .)  # fix any formatting issues
```

### 6. Bump flake.nix version

```nix
version = "0.X.Y";
```

Note: if `go.mod`/`go.sum` changed, the Nix `vendorHash` in `flake.nix` will be stale. Set it to `pkgs.lib.fakeHash`, run `nix build .`, and copy the correct hash from the error.

### 7. Commit

```bash
git add -A && git commit -m "merge: upstream tomasz-tomczyk/crit vX.Y.Z"
```
