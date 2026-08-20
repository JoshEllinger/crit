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
| v0.16.5 | 2026-06-29 | internal/config/config.go (OpenCmd added; kept ShareURL=""), test/e2e/tests/settings-panel.spec.ts (config card assertions) | Massive repo layout refactor (v0.16.3): internal/ packages, web/, integrations/ moved to root; 127 Go files needed module path updated from tomasz-tomczyk to JoshEllinger; isSessionFile exported as IsSessionFile for cross-package use. The "configuration cards" test was re-resolved to match fork rendering (Account/Share cards absent with empty share_url). |
| v0.18.0 | 2026-07-13 | cmd/crit/main.go (live/preview arg-check moved to cli_routing.go — took upstream structure), internal/live/live.go (new `internal/focus` import), internal/server/server.go (import block + handleFiles: kept our IsAbs/IsSessionFile early-return, adopted upstream's tightened `filepath.Clean`+`../` traversal check for the relative-path branch), internal/server/server_test.go (new `internal/vcs` import) | Merged through upstream/main (3 commits past the v0.18.0 tag: story-review hardening, story-mode rendering polish, readme fix — none warranted a separate fork version). Brought in story mode, command hooks, prompt/trust features, live-mode reconnect (new `github.com/gorilla/websocket` dep — vendorHash in flake.nix needs regenerating via CI's `check-nix-build.sh` since no local Nix). ~35 new upstream Go files needed the `tomasz-tomczyk`→`JoshEllinger` import path fix (story, prompt, hooks packages + new test files). settings-panel.spec.ts auto-merged cleanly this time (no touch in this range) — fork's Account/Share divergence comment still accurate, left as-is. Custom changes verified intact: IsSessionFile, handleFiles IsAbs block, no homebrew tap step. |
| v0.18.2 | 2026-07-31 | ~20 files, all pure import-block conflicts from upstream package moves/additions (`notify`, `reviewpath` new packages; several test files gained/dropped helper-package imports) across cmd/crit and internal/{comment,config,focus,github,live,review,server,session,share} — resolved by taking upstream's import list and translating `tomasz-tomczyk`→`JoshEllinger`. Sole logic conflict: internal/session/session.go `GetFileSnapshotFromDisk` — took upstream's tightened traversal check (Clean+`..` reject, then symlink-resolve both the target and repo root and compare resolved paths) over our earlier repoRoot-only symlink resolution, same precedent as v0.18.0's handleFiles hardening. | Merged through the v0.18.2 tag only (upstream had 4 more unreleased commits past it — comment-nav fix, agent-reply-flow fix, marker-CSS fallback, a stylelint bump — deferred to the next merge). ~11 new upstream Go files (cli_dispatch.go, cli_resolve.go, plans.go, plan_cli_main_test.go, and several new *_test.go files) needed the import path fix — none were conflicts since fork had no prior version of those files. go.mod/go.sum ended up byte-identical to upstream's (module path line aside) so flake.nix's `vendorHash` needed no change; `version` auto-merged cleanly to `0.18.2`. settings-panel.spec.ts had no changes in this range — didn't conflict. Custom changes verified intact: IsSessionFile, handleFiles IsAbs block, ShareURL default, no homebrew tap step. Full `go test ./...` (all packages) and `gofmt -l .` clean; E2E left to CI (no local node_modules in the fresh worktree used for this merge) rather than provisioning a new one. |
| v0.18.4 | 2026-08-11 | 9 files, all pure import-block conflicts (GitLab merge-request integration, coverage-workflow tooling, and other new upstream packages added imports) in internal/comment/{daemon_config_precedence_test,flags,list_test,run}.go, internal/github/push_cli_flags_test.go, internal/live/live.go, internal/preview/preview.go, internal/review/review_file_test.go, internal/server/daemon_cli.go — resolved by taking upstream's import list and translating `tomasz-tomczyk`→`JoshEllinger`. No logic conflicts this round: `internal/config/config.go` (ShareURL default), `internal/session/session.go` (IsSessionFile), `internal/server/server.go` (handleFiles IsAbs block), and `test/e2e/tests/settings-panel.spec.ts` all auto-merged cleanly with the fork's customizations intact — verified by grep/read afterward rather than assumed. | Merged through the v0.18.4 tag (upstream had 5 more unreleased commits past it on main — a configurable code-font setting, a mermaid bump, a markdown-it major bump, a scroll-position fix, and a share-transport routing fix — deferred to the next merge). go.mod/go.sum ended up byte-identical to the fork's pre-merge state (no dependency changes in this range), so `vendorHash` needed no change; `flake.nix` `version` auto-merged cleanly to `0.18.4`. `integrations/claude-code/.claude-plugin/plugin.json` auto-merged its version bump (1.8.2→1.8.6) along with all the `integrations/*/skills/*.md` content; `install.sh` (fork-local, not present upstream) was untouched. Ran `go generate ./...` afterward to regenerate `cmd/crit/integration_hashes_gen.go` against the merged skill content — came out byte-identical, no drift. **Gotcha this round:** a naive awk-based conflict resolver (strip `<<<<<<< HEAD`/`=======`/`>>>>>>>` markers and keep both sides) produced duplicate import lines instead of upstream-only imports, since HEAD's and upstream's import sets overlapped almost entirely — caught immediately by `go build` failing with "redeclared"/"imported and not used", fixed with a targeted per-file de-duplication pass scoped to each `import (...)` block. Full `go build`, `go test ./...` (all packages), `go vet ./...`, and `gofmt -l .` clean. Confirmed no unrelated files touched (`git status --porcelain` had zero untracked entries post-merge). E2E left to CI per the v0.18.2 precedent. Pre-existing (not touched by this merge, not a regression): `internal/server/server.go`'s `CheckForUpdates` still points at `tomasz-tomczyk/crit`'s GitHub releases for the update-notification check rather than the fork's own releases — flagged to the user, left as-is pending a decision on whether the fork wants that behavior changed. |
| v0.19.0 | 2026-08-20 | 11 files, all pure import-block conflicts from upstream's new `internal/forge`/`internal/gitlab` abstraction (added to support the GitLab MR integration merged at v0.18.4): cmd/crit/{cli_handlers,wire,wire_test}.go, internal/github/{github,github_test}.go, internal/server/{handle_share_round_test,picker,picker_http_test,picker_stack_test}.go, internal/session/{session,types}.go — resolved by taking upstream's import list (a strict superset of the fork's in every case) and translating `tomasz-tomczyk`→`JoshEllinger`. No logic conflicts. Custom changes verified intact: IsSessionFile, handleFiles IsAbs block, ShareURL default, no homebrew tap step, settings-panel.spec.ts assertions (auto-merged cleanly again). | Merged through the v0.19.0 tag (21 upstream commits since v0.18.4): the deferred configurable-code-font/mermaid-bump/markdown-it-15 items from the v0.18.4 note landed in this range, plus native markdown table rendering, GitHub-compatible HTML in comment markdown, and several fixes. go.mod gained two new transitive deps (`github.com/go-text/typesetting`, `golang.org/x/image`, pulled in by the forge/gitlab work) — left `flake.nix`'s `vendorHash` untouched and let CI's `nix build` validate it rather than guessing without local Nix; it passed unchanged. `flake.nix` `version` auto-merged cleanly to `0.19.0`. Ran `go generate ./...` for `integration_hashes_gen.go` — no drift. **Gotcha this round:** CI's `lint` job failed on a file untouched by this merge (`internal/live/proxy.go`) — golangci-lint's `latest` had picked up a newer staticcheck rule flagging `httputil.ReverseProxy.Director` as deprecated since Go 1.26. Confirmed via history (fork's v0.18.4 PR lint was clean 9 days prior) this was linter drift, not a merge regression; fixed by migrating to `Rewrite` (behavior-preserving, verified against the existing `proxy_test.go` suite). **Second gotcha:** `e2e-windows` failed non-deterministically-looking-deterministic (all 14 failures confined to one spec file, both attempts) on the `windows-2025-vs2026` preview runner image; before assuming a real regression, confirmed upstream's own PR adding the same new tests (tomasz-tomczyk/crit#834) passed `e2e-windows` cleanly on their CI, and a bare re-run of the job here passed clean too — treated as preview-runner flakiness, not a code issue. Full `go build`, `go test ./...` (one `TestDaemonLifecycle` flake under full-suite load, confirmed passing 3/3 in isolation), `go vet ./...`, and `gofmt -l .` clean. |

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
- **test/e2e/tests/settings-panel.spec.ts** — The "settings pane shows configuration cards" test WILL conflict on every merge. Upstream asserts that "Account" and "Sharing enabled" cards are visible because upstream defaults `share_url` to `"https://crit.md"`. Our fork defaults `share_url` to `""`, so those cards do NOT render in git-mode. After merging, resolve this test toward the fork's actual rendering: keep "Agent Command" and the `AI Integration|Integration Available` assertions; assert "Account" has count 0 (not visible); assert "Share" (not "Sharing enabled") is visible. Do NOT blindly take upstream's Account/Share assertions.

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
