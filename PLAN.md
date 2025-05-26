# Implementation Plan

## Core Behavior

1. **Unified View**: Show all labels from both repos in one list
   - Target-only labels appear first with "(dest only)" marker
   - Visual separator between target-only and source labels
   - Source labels are pre-selected by default

2. **Label Matching**: Labels are matched by name (case-insensitive)
   - If a label exists in both repos, use the source version
   - Never create duplicates

3. **User Actions**:
   - Toggle any label on/off
   - Unchecking a target-only label = delete it
   - Checking a source label = copy it
   - Final confirmation shows exact end state

## Non-Functional Requirements

### Logging
- Use stderr for all logs to keep stdout clean
- Log levels: DEBUG (hidden by default), INFO, ERROR
- Enable debug with `-d` or `--debug` flag
- Log format: `[LEVEL] message`
- Key events to log:
  - GitHub CLI detection
  - API calls made
  - Label operations (create/update/delete)
  - Errors with full context

### Error Handling
1. **GitHub CLI not found**:
   ```
   Error: GitHub CLI (gh) is required but not found.
   Install it from: https://cli.github.com
   ```

2. **Not authenticated**:
   ```
   Error: Not authenticated with GitHub.
   Run: gh auth login
   ```

3. **No permissions**:
   ```
   Error: You don't have permission to manage labels in owner/repo.
   Ensure you have write access to the repository.
   ```

4. **Network/API errors**: Show GitHub's error message directly

### Terminal Compatibility
- Show colored blocks before hex codes: `█ #d73a4a bug`
- Detect color support:
  - 24-bit color: Full color blocks
  - 256 colors: Approximate to nearest color
  - 16 colors: Show colored text instead of blocks
  - No color: Show hex codes only
- Use `github.com/mattn/go-isatty` to detect TTY
- Gracefully degrade in non-interactive environments

## Implementation Details

### CLI Structure
```
main.go         - cobra setup, command handling
labels.go       - fetch/create/update/delete via gh api
picker.go       - promptui interactive selection
display.go      - color rendering, terminal output
types.go        - Label struct and related types
logger.go       - Logging utilities
```

### GitHub API Calls
```bash
# Fetch labels
gh api repos/{owner}/{repo}/labels --paginate

# Create label
gh api repos/{owner}/{repo}/labels -f name=X -f color=X -f description=X

# Update label
gh api repos/{owner}/{repo}/labels/{name} --method PATCH -f color=X -f description=X

# Delete label
gh api repos/{owner}/{repo}/labels/{name} --method DELETE
```

### Edge Cases

1. **No labels in source**: Show error "No labels found in {owner}/{repo}"
2. **Empty selection**: Show warning "No labels selected. Nothing to do."
3. **Ctrl+C handling**: Clean exit, no partial operations
4. **Large label sets**: Handle pagination properly

### Testing

1. Mock `gh` commands for unit tests
2. Test picker state management
3. Test confirmation logic
4. Test error handling for all edge cases
5. Test color detection and fallback

### Display Format

Picker display with full color support:
```
  [✓] █ bug                 #d73a4a
  [✓] █ documentation       #0075ca
  [ ] █ duplicate           #cfd3d7
```

Limited color support:
```
  [✓] bug                 #d73a4a [red]
  [✓] documentation       #0075ca [blue]
  [ ] duplicate           #cfd3d7 [gray]
```

Final confirmation:
```
Final state for myorg/myproject:
  ✓ █ bug                #d73a4a
  ✓ █ documentation      #0075ca
  ✓ █ enhancement        #a2eeef

Actions:
  • Create 2 labels
  • Delete 1 label (duplicate)
  • Keep 1 existing label

Proceed? [y/N]:
```