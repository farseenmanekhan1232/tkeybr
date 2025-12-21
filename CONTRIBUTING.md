# Contributing to Hyperkey

Thank you for your interest in contributing to Hyperkey! 🎉

## How to Contribute

### Reporting Bugs
- Open an issue with a clear description
- Include steps to reproduce
- Mention your OS and terminal emulator

### Suggesting Features
- Open an issue with the `enhancement` label
- Describe the use case and expected behavior

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Build: `go build ./...`
6. Commit with clear messages
7. Push and open a PR

### Code Style
- Follow standard Go conventions (`gofmt`)
- Keep functions small and focused
- Add comments for complex logic

### Areas for Contribution

#### High Priority
- [ ] More keyboard layouts (Workman, Colemak-DH, etc.)
- [ ] Cloud sync with Google account
- [ ] Settings UI with interactive controls
- [ ] Unit tests

#### Nice to Have
- [ ] Typing games/challenges
- [ ] Multiplayer races
- [ ] Custom text/lessons
- [ ] More languages (non-English phonetic patterns)

## Development Setup

```bash
# Clone
git clone https://github.com/yourusername/hyperkey.git
cd hyperkey

# Install dependencies
go mod download

# Build
go build -o hyperkey ./cmd/hyperkey

# Run
./hyperkey
```

## Project Structure

```
hyperkey/
├── cmd/hyperkey/          # Entry point
├── internal/
│   ├── app/               # BubbleTea application
│   ├── engine/            # Typing engine, stats, layouts
│   ├── auth/              # Google OAuth
│   ├── storage/           # Local persistence
│   └── ui/
│       ├── components/    # UI components
│       └── styles/        # Lipgloss styles
└── README.md
```

## Questions?

Open an issue or reach out! We're happy to help.
