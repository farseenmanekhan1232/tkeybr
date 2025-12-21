# tkeybr

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terminal-based keyboard practice application inspired by [keybr.com](https://keybr.com). Learn touch typing with adaptive learning that focuses on your weak keys.

## Features

### Adaptive Learning (keybr-style)
- **Progressive Letter Unlocking** - Start with home row keys, unlock more as you improve
- **Confidence-Based Progression** - Each letter tracked individually (accuracy + speed)
- **Weak Key Focus** - Struggling letters appear more frequently
- **Phonetic Words** - Pronounceable pseudo-words for natural typing

### Multiple Keyboard Layouts
| Layout | Home Row | Same Finger |
|--------|----------|-------------|
| QWERTY | 30% | 19% |
| Dvorak | 68% | 18% |
| Colemak | 70% | 15% |

### Real-time Statistics
- Words Per Minute (WPM) - correct characters only
- Accuracy percentage
- Keyboard heatmap
- Per-key error tracking

## Installation

### Homebrew (macOS/Linux)
```bash
brew tap farseenmanekhan1232/tap
brew install tkeybr
```

### Linux Packages
Download `.deb`, `.rpm`, or `.apk` from [Releases](https://github.com/farseenmanekhan1232/tkeybr/releases).

```bash
# Debian/Ubuntu
sudo dpkg -i tkeybr_*.deb

# Fedora/RHEL
sudo rpm -i tkeybr_*.rpm

# Alpine
sudo apk add --allow-untrusted tkeybr_*.apk
```

### Windows
Download the `.zip` from [Releases](https://github.com/farseenmanekhan1232/tkeybr/releases) and extract.

### Go Install
```bash
go install github.com/farseenmanekhan1232/tkeybr/cmd/tkeybr@latest
```

### From Source
```bash
git clone https://github.com/farseenmanekhan1232/tkeybr.git
cd tkeybr
go build -o tkeybr ./cmd/tkeybr
./tkeybr
```

## Usage

```bash
./tkeybr
```

### Controls
| Key | Action |
|-----|--------|
| Up/Down or j/k | Navigate menu |
| Enter | Select |
| Esc | Return to menu |
| q or Ctrl+C | Quit |

### How It Works
1. **Start with 4 letters** (f, j, d, k on QWERTY)
2. **Type practice words** - errors are recorded, no backspace allowed
3. **Build confidence** - each letter needs 70% confidence to unlock the next
4. **Progress automatically** - new letters unlock as you master current ones

## Project Structure

```
tkeybr/
├── cmd/
│   ├── tkeybr/          # CLI entry point
│   └── server/            # SSH server (for tkeybr.io)
├── internal/
│   ├── app/               # BubbleTea application
│   ├── engine/            # Typing engine, stats, layouts
│   ├── auth/              # Google OAuth
│   ├── storage/           # Local persistence (~/.tkeybr/)
│   └── ui/
│       ├── components/    # Keyboard, typing area, stats, charts
│       └── styles/        # Lipgloss color theme
├── .goreleaser.yaml       # Cross-platform release config
└── .github/workflows/     # CI/CD
```

## Tech Stack

- **[Go](https://go.dev)** - Fast, simple, compiles to single binary
- **[BubbleTea](https://github.com/charmbracelet/bubbletea)** - TUI framework (Elm architecture)
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[Wish](https://github.com/charmbracelet/wish)** - SSH server for TUI apps

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Coming Soon

- **SSH Access** - `ssh tkeybr.io` (no installation needed!)
- **Winget** - `winget install tkeybr` for Windows
- More layouts (Workman, Colemak-DH)
- Cloud sync with Google account
- Interactive settings UI
- Typing games
- Tests

## License

MIT License - see [LICENSE](LICENSE)

## Acknowledgments

- Inspired by [keybr.com](https://keybr.com)
- Built with [Charm](https://charm.sh) tools
