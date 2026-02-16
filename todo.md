# PlanReview — TODO

## Features (From Spec, Not Yet Implemented)

- [ ] **Comment collapse/expand**: Comments can be collapsed/expanded (spec section "Displayed comments"). Currently all comments are always expanded.
- [ ] **Empty/binary file handling**: Show user-friendly messages for empty files or non-markdown files (spec "Edge Cases").
- [ ] **Very large files**: Test with files up to ~10k lines, ensure no performance issues.

## UI Refinements

- [ ] **Mobile/responsive**: Basic responsive CSS exists but untested on small screens.

## Pre-publish (Done)

- [x] **Go tooling**: go.mod updated to 1.26, golangci-lint clean, gofmt clean
- [x] **Security review**: Fixed path traversal in `/files/`, added request body limits (1MB), added HTTP server timeouts
- [x] **JS review**: Fixed filename XSS in innerHTML; `html:false` on comment renderer already safe; `html:true` on doc renderer is intentional (local tool)
- [x] **Publish readiness**: LICENSE file added, README rewritten. `YOUR_USERNAME` placeholders and go.mod module path TBD when GitHub repo created.
- [x] **Unit tests**: 37 tests across server_test.go, document_test.go, output_test.go — covers API handlers, CRUD, validation, path traversal, output generation, concurrent access

## Future Enhancements (Post-v1)

- [ ] **GitHub Actions release workflow**: Cross-compile binaries on tagged releases (spec has a workflow sketch).
- [ ] **Homebrew tap**: `brew install planreview`.
- [ ] **Comment resolution**: Mark comments as "resolved" (like GitHub), visually collapsed but still in review file.
- [ ] **Diff view**: After the AI agent updates the plan, show what changed alongside original comments.
- [ ] **Multiple reviewers**: Support `--author "Tom"` for team review scenarios.
- [ ] **Configurable comment format**: Let users customize the review comment prefix in output.
- [ ] **Export formats**: Export comments as GitHub Issues, Linear tickets, or TODO list.
