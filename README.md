# Gable

Safely copy GitHub labels between repositories. Review everything before making changes - you're always in control.

## Install

Download the latest binary:
- [macOS (Intel)](https://github.com/robert-claypool/gable/releases/latest/download/gable-darwin-amd64)
- [macOS (Apple Silicon)](https://github.com/robert-claypool/gable/releases/latest/download/gable-darwin-arm64)
- [Linux](https://github.com/robert-claypool/gable/releases/latest/download/gable-linux-amd64)
- [Windows](https://github.com/robert-claypool/gable/releases/latest/download/gable-windows-amd64.exe)

Make it executable and move to your PATH:
```bash
chmod +x gable
mv gable /usr/local/bin/
```

## Usage

```bash
gable -h # show help
gable owner/source owner/dest
```

This opens an interactive picker showing all labels from both repos (source and destination). Use arrow keys to navigate, Space to toggle, Enter to confirm selections.

```
Current state → Desired state for myorg/myproject:

  [✓] NeedsFix            #aa0000  (dest only)
  ────────────────────────────────────────────────
  [✓] bug                 #d73a4a
  [✓] documentation       #0075ca
  [ ] duplicate           #cfd3d7
  [✓] enhancement         #a2eeef

  Space: toggle  ↑/↓: navigate  Enter: confirm  q: quit
```

## Requirements

**GitHub CLI is required.** Gable uses the GitHub CLI to interact with GitHub.

1. [Install GitHub CLI](https://cli.github.com)
2. Run `gh auth login` to authenticate
3. Ensure you have permission to manage labels in the destination repository

## Options

- `-v, --verbose` - Show label descriptions
- `-d, --debug` - Show debug logs
- `-h, --help` - Show help

## License

MIT
