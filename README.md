# PlanReview

A lightweight CLI tool for reviewing markdown plans with GitHub PR-style inline comments. Built for the workflow of iterating on plans with AI coding agents.

You write a plan in markdown. You run `planreview` against it. You leave inline comments — single-line or multi-line ranges, just like a GitHub PR review. The tool writes a `.review.md` file in real-time with your comments interleaved, ready to hand back to your AI agent at any moment.

## Install

### From GitHub Releases (recommended)

Download the latest binary for your platform from [Releases](https://github.com/YOUR_USERNAME/planreview/releases):

| Platform           | Binary                          |
|--------------------|---------------------------------|
| macOS (Apple Silicon) | `planreview-darwin-arm64`     |
| macOS (Intel)      | `planreview-darwin-amd64`       |
| Linux (x86_64)     | `planreview-linux-amd64`        |
| Linux (ARM64)      | `planreview-linux-arm64`        |

On macOS, after downloading:

```bash
# Move to a directory on your PATH
mv ~/Downloads/planreview-darwin-arm64 /usr/local/bin/planreview

# Make it executable
chmod +x /usr/local/bin/planreview

# On first run, macOS will block it. Either:
# Option A: Right-click → Open in Finder, then confirm
# Option B: Remove the quarantine attribute:
xattr -d com.apple.quarantine /usr/local/bin/planreview
```

### Build from Source

Requires Go 1.22+ (install via [asdf](https://asdf-vm.com/), Homebrew, or [go.dev](https://go.dev/dl/)):

```bash
# With asdf
asdf plugin add golang
asdf install golang latest
asdf global golang latest

# Clone and build
git clone https://github.com/YOUR_USERNAME/planreview.git
cd planreview
go build -o planreview .

# Optionally move to your PATH
mv planreview /usr/local/bin/
```

## Usage

```bash
# Review a markdown file
planreview plan.md
# → Opens browser at http://localhost:<port>
# → Leave inline comments, GitHub PR-style
# → plan.review.md is written in real-time

# Specify a port
planreview --port 3000 plan.md

# Don't auto-open browser
planreview --no-open plan.md
```

## Workflow

```bash
# 1. AI agent generates a plan
agent "Write a plan for the new auth service" > auth-plan.md

# 2. Review it
planreview auth-plan.md
# → Leave comments in the browser
# → On finish, the prompt is copied to your clipboard

# 3. Hand the review file back to your agent
agent "I've left review comments in auth-plan.review.md — please address
       each comment and update the plan accordingly."

# 4. Review the updated plan
planreview auth-plan.md
# → Fresh session, repeat as needed
```

### Output Files

| File | Purpose |
|------|---------|
| `plan.review.md` | Your plan with comments interleaved — hand this to your AI agent |
| `.plan.comments.json` | Hidden file for resuming sessions (ignore this) |

## How It Works

- Comments reference **source line numbers** in the original `.md` file
- The `.review.md` is the original plan with comments inserted as block quotes
- If you reopen a file with existing comments (and the file hasn't changed), your previous comments are restored
- The server exits when you click "Finish Review" or press Ctrl+C

## License

MIT
