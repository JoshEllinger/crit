# PlanReview — Development Guide

## What This Is

A single-binary Go CLI tool that opens a browser-based UI for reviewing markdown files with GitHub PR-style inline commenting. Comments are written in real-time to a `.review.md` file designed to be handed to Claude.

## Project Structure

```
planreview/
├── main.go              # Entry point: CLI parsing, server setup, graceful shutdown, browser open
├── server.go            # HTTP handlers: REST API (comments CRUD, document, finish, stale)
├── document.go          # Core state: file loading, comment storage, JSON/review file persistence
├── output.go            # Generates .review.md (original markdown + interleaved comment blockquotes)
├── frontend/
│   ├── index.html       # Complete SPA — all HTML, CSS, JS in one file (~1500 lines)
│   ├── markdown-it.min.js    # Markdown parser (provides source line mappings via token.map)
│   ├── highlight.min.js      # Syntax highlighter core
│   ├── hljs-*.min.js         # Language packs (js, ts, go, python, elixir, etc.)
│   └── hljs-{dark,light}.css # Highlight themes (embedded inline in index.html via data-theme scoping)
├── go.mod
├── Makefile             # build / build-all (cross-compile)
├── test-plan.md         # Sample file for development testing
└── README.md
```

## Key Architecture Decisions

1. **All frontend assets embedded** via Go's `embed.FS` — produces a true single binary
2. **No frontend build step** — vanilla JS, no npm/webpack/framework
3. **markdown-it for parsing** — chosen because it provides `token.map` (source line mappings per block)
4. **Block-level splitting** — lists, code blocks, tables, blockquotes are split into per-item/per-line/per-row blocks so each source line is independently commentable
5. **Comments reference source line numbers** — the `.review.md` output uses `> **[REVIEW COMMENT — Lines X-Y]**:` format
6. **Real-time output** — `.review.md` and `.comments.json` written on every comment change (200ms debounce)
7. **GitHub-style gutter interaction** — click-and-drag on line numbers to select ranges

## Build & Run

```bash
go build -o planreview .          # Build
./planreview test-plan.md         # Run (opens browser)
./planreview --no-open --port 3000 test-plan.md  # Headless on fixed port
```

## API Endpoints

- `GET  /api/document` — raw markdown content + filename
- `GET  /api/comments` — all comments
- `POST /api/comments` — add comment `{start_line, end_line, body}`
- `PUT  /api/comments/:id` — edit comment `{body}`
- `DELETE /api/comments/:id` — delete comment
- `POST /api/finish` — write final files and shut down server
- `GET  /api/stale` — check if file changed since last session
- `DELETE /api/stale` — dismiss stale notice

## Frontend Architecture (index.html)

The trickiest part is **source line mapping**. The approach:

1. Parse markdown with `markdown-it` to get tokens with `token.map` (source line ranges)
2. `buildLineBlocks()` walks the token stream and creates a flat array of commentable blocks
3. Container tokens (lists, tables, blockquotes) are drilled into — each list item, table row, or blockquote child becomes its own block
4. Code blocks (`fence` tokens) are split into per-line blocks with syntax highlighting preserved via `splitHighlightedCode()` which handles `<span>` tags crossing line boundaries
5. Each block gets a gutter entry with its source line number(s)
6. Comments are keyed by `end_line` and displayed after their referenced block

### Known Complexities

- **markdown-it token.map quirks**: The last item in a list often claims a trailing blank line (e.g., map `[94, 96]` for a single-line item). The code trims trailing blank lines from item ranges.
- **Table separator lines** (`|---|---|`): Not represented in tokens, appear as gap lines. Detected via regex and hidden with CSS.
- **Per-row tables**: Each row wrapped in its own `<table>` with `table-layout: fixed` + `<colgroup>` for column alignment.
- **Highlighted code splitting**: `splitHighlightedCode()` tracks open `<span>` tags across lines to properly close/reopen them.

## Output Files

| File | Description |
|------|-------------|
| `plan.review.md` | Original markdown + comments as blockquotes — hand to Claude |
| `.plan.comments.json` | Hidden dotfile for resume support (stores file hash) |
