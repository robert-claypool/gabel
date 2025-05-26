# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gabel is a CLI tool that helps users copy GitHub labels between repositories with an interactive picker.

## Commands

```bash
# Build
go build -o gabel

# Run
./gabel owner/source owner/dest

# Test
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Architecture

Simple Go CLI that wraps GitHub CLI (`gh`):
- `main.go` - CLI entry point using cobra
- `labels.go` - GitHub API calls via `gh api`
- `picker.go` - Interactive selection using promptui
- `display.go` - Terminal colors and formatting

## Key Implementation Notes

- Uses `gh api` for all GitHub operations - no API tokens needed
- Interactive picker shows unified view of labels from both repos
- Never overwrites without explicit user confirmation
- Minimal dependencies: cobra for CLI, promptui for picker, fatih/color for colors

## Development Guidelines

### Core Principles

1. **Fail Fast, Fail Clearly**: Never swallow errors. Let processes die with actionable error messages
2. **Observable by Default**: User-facing logs for progress, DEBUG for developer details
3. **Test-Driven Development**: Write tests first, mock external dependencies
4. **Composable Unix Philosophy**: Support pipes, meaningful exit codes, streamable output
5. **Strict Meaningful Types**: Use type safety to your advantage
6. **Generic Infrastructure**: Keep non-functional code domain-agnostic. Use generic examples (foo/bar) in tests, utilities, and infrastructure

### Code Style

- Focus comments on WHY not WHAT. Document architectural decisions
- Error messages: include a suggestion to fix the problem if possible
- Use golangci-lint

### Logging

- **DEBUG** - Development details, hidden by default (enabled with -d flag)
- **INFO** - User progress updates (printed to stdout)
- **ERROR** - Problems that prevent operation (printed to stderr)

Note: Go doesn't have built-in log levels like WARN. We use fmt.Printf for INFO and fmt.Fprintf(os.Stderr) for ERROR.

### Testing Requirements

- 80% code coverage target

**Coverage Pragmatism**: The 80% target balances thoroughness with practicality. Integration tests provide additional confidence beyond unit test coverage. Focus testing on business logic and error paths.

## Console Output Guidelines

- **No emojis** in console output. Use plain ASCII or common Unicode when necessary
- Use color for emphasis and clarity
- For status indicators, use text like `[OK]`, `[ERROR]`, `>>` instead of emoji

## Writing Style

Prefer simple, common words:
- "need" not "necessary"
- "use" not "utilize"
- "help" not "assist"
- "start" not "commence"

## Comment Standards

Follow these principles:

1. Concise and Direct: Begin each comment with the main point or rationale
2. Avoid Jargon Unless Essential: Use acronyms or industry-specific terms only when essential
3. Document Lessons Learned: Include reasoning behind decisions
4. Colocate Comments: Place comments close to the relevant code
5. Focus on Why, Not What: Explain rationale rather than stating the obvious
6. Highlight Non-Obvious Details: Clarify units, contexts, or gotchas
7. Warn About Risks: Clearly highlight critical sections
8. No Emojis in Code: Use ASCII. Prefer "TODO", "FIXME", "[OK]", "[FAIL]" over emoji