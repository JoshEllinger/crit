# Upstream Merge Plan

## Context

We maintain a fork of `tomasz-tomczyk/crit` with three custom changes that need to survive each upstream merge:

1. **Files outside repo root fix** — When crit is given a file outside the repo root (e.g. `~/.claude/plans/foo.md`), `filepath.Rel` produces a `../..` path that the `/files/` endpoint rejects. Our fix: use absolute paths and validate against known session files via `isSessionFile()`. Upstream does NOT have this fix.
2. **No Homebrew tap push** — upstream's release.yml pushes to `tomasz-tomczyk/homebrew-tap` using a `HOMEBREW_TAP_TOKEN` secret we don't have. Strip this step after merge if present.
3. **`--global` flag for `crit install`** — `crit install --global claude-code` writes to `~/.claude/commands/crit.md` for user-wide availability. Upstream does not have this.

Items that are NOT concerns:
- The Nix flake and marketplace JSON files reference `tomasz-tomczyk/crit` but are inert in our fork.
- `install.sh` — our cross-platform install script, local-only, just keep it.

## Merge History

| Version | Date | Conflicts | Notes |
|---------|------|-----------|-------|
| v0.8.3 | 2026-04-10 | main.go, session.go, server.go | Initial merge from v0.7.0 |
| v0.9.2 | 2026-04-17 | main.go (help text), session.go (isSessionFile) | Clean merge, 20 upstream commits |

## What's new in v0.9.2 (upstream additions)

- Content-based comment anchoring for carry-forward positioning
- Tab-ready title indicator for background review rounds
- Opportunistic background cleanup of stale reviews and sessions
- Orphaned comments surfaced on removed files
- `crit auth` subcommands (login/logout/whoami)
- `crit status` and `crit cleanup` subcommands
- Only poll git status while waiting for agent edits
- Hide share button in git mode
- Various fixes: comment index perf, CSS, test isolation, session key dedup, debounced writes, unified-diff selection, hljs markdown fences, git --no-optional-locks, crit comment --reply-to subdirectory fix
- Quality guardrails (linting, accessibility, rules)
- Favicons and web app manifest
- "review file" terminology replaces ".crit.json" references
- Prevent agents from proactively resolving review comments

## Steps for future merges

### 1. Merge upstream/main (or tag)

```bash
git fetch upstream --tags
git merge v0.X.Y
```

### 2. Resolve conflicts

Expected conflict areas (our custom code vs upstream changes):

- **main.go** — `printHelp()` has our `--global` flag text. Take upstream additions, keep `[--global]` on the install line.
- **session.go** — Our `isSessionFile()` method. Keep it; upstream doesn't have it.
- **server.go** — Our `filepath.IsAbs` block in `handleFiles`. Keep it; upstream doesn't have it.

### 3. Verify custom changes survived

```bash
# Files-outside-root fix
grep -n 'isSessionFile' session.go server.go
grep -n 'filepath.IsAbs(reqPath)' server.go

# --global flag
grep -n '\-\-global' main.go
```

### 4. Sync local skill

Compare `.claude/commands/crit.md` with upstream's `integrations/claude-code/skills/crit/SKILL.md` and sync content changes while keeping local frontmatter format.

### 5. Build, test, install

```bash
go test ./...
go build -o ~/.local/bin/crit .
```

### 6. Commit

```bash
git add -A && git commit -m "merge: upstream tomasz-tomczyk/crit vX.Y.Z"
```
