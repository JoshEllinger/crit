# Crit Integrations

Drop-in configuration files that teach your AI coding tool to write plans, launch Crit for review, and wait for your feedback before implementing.

Copy the file or files for your tool into your project.

| Tool | File to copy | Destination in your project |
|------|-------------|----------------------------|
| Claude Code | `claude-code/crit.md` | `.claude/commands/crit.md` |
| Claude Code | `claude-code/crit-comment.md` | `.claude/commands/crit-comment.md` |
| Claude Code | `claude-code/CLAUDE.md` (optional) | Append to your `CLAUDE.md` |
| Cursor | `cursor/crit-command.md` | `.cursor/commands/crit.md` |
| Cursor | `cursor/crit-comment.md` | `.cursor/commands/crit-comment.md` |
| Cursor | `cursor/crit.mdc` (optional) | `.cursor/rules/crit.mdc` |
| OpenCode | `opencode/crit.md` | `.opencode/commands/crit.md` |
| OpenCode | `opencode/crit-comment.md` | `.opencode/commands/crit-comment.md` |
| OpenCode | `opencode/SKILL.md` | `.opencode/skills/crit-review/SKILL.md` |
| Windsurf | `windsurf/crit.md` | `.windsurf/rules/crit.md` |
| Windsurf | `windsurf/crit-comment.md` | `.windsurf/rules/crit-comment.md` |
| GitHub Copilot | `github-copilot/crit.prompt.md` | `.github/prompts/crit.prompt.md` |
| GitHub Copilot | `github-copilot/crit-comment.md` | `.github/prompts/crit-comment.prompt.md` |
| GitHub Copilot | `github-copilot/copilot-instructions.md` (optional) | Append to `.github/copilot-instructions.md` |
| Aider | `aider/CONVENTIONS.md` | Append to your `CONVENTIONS.md` |
| Aider | `aider/crit-comment.md` | Copy to your project root |
| Cline | `cline/crit.md` | `.clinerules/crit.md` |
| Cline | `cline/crit-comment.md` | `.clinerules/crit-comment.md` |

## What these do

All integrations follow the same pattern:

1. **Plan first** - the agent writes an implementation plan as a markdown file before writing any code
2. **Launch Crit** - the agent runs `crit $PLAN_FILE` to open the plan for review in your browser
3. **Address feedback** - after review, the agent reads `.crit.json` to find your inline comments and revises the plan
4. **Implement after approval** - only after you approve does the agent write code

Claude Code, Cursor, OpenCode, and GitHub Copilot all support a `/crit` slash command that automates the full loop: find the plan, launch Crit, read comments, revise, and signal for another round. OpenCode can also load the `crit-review` skill on demand.

Each integration also includes a `crit-comment` skill that teaches your agent to use `crit comment` to add inline review comments programmatically — no browser needed. The agent learns the syntax and can leave comments on specific lines or ranges as part of its workflow.
