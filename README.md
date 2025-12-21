# Hyperkey ⌨️

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terminal-based keyboard practice application inspired by [keybr.com](https://keybr.com). Learn touch typing with adaptive learning that focuses on your weak keys.

![Hyperkey Demo](https://via.placeholder.com/800x400?text=Demo+Coming+Soon)

## ✨ Features

### 🎯 Adaptive Learning (keybr-style)
- **Progressive Letter Unlocking** - Start with home row keys, unlock more as you improve
- **Confidence-Based Progression** - Each letter tracked individually (accuracy + speed)
- **Weak Key Focus** - Struggling letters appear more frequently
- **Phonetic Words** - Pronounceable pseudo-words for natural typing

### ⌨️ Multiple Keyboard Layouts
| Layout | Home Row | Same Finger |
|--------|----------|-------------|
| QWERTY | 30% | 19% |
| Dvorak | 68% | 18% |
| Colemak | 70% | 15% |

### 📊 Real-time Statistics
- Words Per Minute (WPM) - correct characters only
- Accuracy percentage
- Keyboard heatmap
- Per-key error tracking

### 🔐 Google Account (Optional)
- Sign in with Google (device code flow)
- Local progress saved to `~/.hyperkey/`

## 🚀 Installation

### From Source
```bash
git clone https://github.com/farseen/hyperkey.git
cd hyperkey
go build -o hyperkey ./cmd/hyperkey
./hyperkey
```

### Go Install
```bash
go install github.com/farseen/hyperkey/cmd/hyperkey@latest
```

## 🎮 Usage

```bash
./hyperkey
```

### Controls
| Key | Action |
|-----|--------|
| ↑/↓ or j/k | Navigate menu |
| Enter | Select |
| Esc | Return to menu |
| q or Ctrl+C | Quit |

### How It Works
1. **Start with 4 letters** (f, j, d, k on QWERTY)
2. **Type practice words** - errors are recorded, no backspace allowed
3. **Build confidence** - each letter needs 70% confidence to unlock the next
4. **Progress automatically** - new letters unlock as you master current ones

## 📁 Project Structure

```
hyperkey/
├── cmd/hyperkey/          # Entry point
├── internal/
│   ├── app/               # BubbleTea application
│   ├── engine/            # Typing engine, stats, layouts
│   ├── auth/              # Google OAuth
│   ├── storage/           # Local persistence (~/.hyperkey/)
│   └── ui/
│       ├── components/    # Keyboard, typing area, stats, charts
│       └── styles/        # Lipgloss color theme
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

## 🛠️ Tech Stack

- **[Go](https://go.dev)** - Fast, simple, compiles to single binary
- **[BubbleTea](https://github.com/charmbracelet/bubbletea)** - TUI framework (Elm architecture)
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - UI components

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Roadmap
- [ ] More layouts (Workman, Colemak-DH)
- [ ] Cloud sync
- [ ] Interactive settings UI
- [ ] Typing games
- [ ] Tests

## 📄 License

MIT License - see [LICENSE](LICENSE)

## 🙏 Acknowledgments

- Inspired by [keybr.com](https://keybr.com)
- Built with [Charm](https://charm.sh) tools
