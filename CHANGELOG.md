# Changelog

All notable changes to Gabel will be documented in this file.

## [1.0.0] - 2025-01-25

### Features
- Interactive label picker with checkbox selection
- Unified view of labels from both source and destination repos
- Smart label matching (case-insensitive)
- Color preview for each label
- Batch select/deselect with 'a' key
- Progress indicators during operations
- Debug mode with `-d` flag
- Verbose mode with `-v` flag to show descriptions
- Clear warnings when labels will be deleted
- Comprehensive error handling with actionable messages

### Technical
- Built with Go and Cobra CLI framework
- Uses GitHub CLI (`gh`) for all API operations
- Terminal color support with automatic detection
- Input validation for label data
- URL encoding for labels with special characters
- 80%+ test coverage on critical paths

### Requirements
- Go 1.20 or higher to build
- GitHub CLI (`gh`) installed and authenticated