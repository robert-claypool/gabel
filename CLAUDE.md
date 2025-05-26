# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gable is a CLI tool that helps users copy GitHub labels between repositories with an interactive picker.

## Commands

```bash
# Build
go build -o gable

# Run
./gable owner/source owner/dest

# Test
go test ./...
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