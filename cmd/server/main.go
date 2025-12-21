package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	gossh "golang.org/x/crypto/ssh"
	"github.com/farseenmanekhan1232/hyperkey/internal/engine"
	"github.com/farseenmanekhan1232/hyperkey/internal/ui/components"
	"github.com/farseenmanekhan1232/hyperkey/internal/ui/styles"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/hyperkey_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

// teaHandler creates a new BubbleTea program for each SSH session
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Get user info from SSH session
	user := s.User()
	pubKey := s.PublicKey()
	
	var userID string
	if pubKey != nil {
		// Use SSH key fingerprint as unique user ID
		userID = gossh.FingerprintSHA256(pubKey)
	} else {
		userID = "guest-" + s.RemoteAddr().String()
	}

	// Get terminal size
	pty, _, active := s.Pty()
	if !active {
		log.Warn("No PTY requested", "user", user)
		return nil, nil
	}

	// Create model with user session info
	m := NewSSHModel(userID, user, pty.Window.Width, pty.Window.Height)

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// SSHModel is the model for SSH sessions (simplified version of full app)
type SSHModel struct {
	userID    string
	username  string
	width     int
	height    int
	view      int // 0=welcome, 1=practice
	
	// Practice state
	typingState     *engine.TypingState
	progressTracker *engine.ProgressTracker
	wordGenerator   *engine.AdaptiveWordGenerator
	keyboard        *components.Keyboard
	lastCharTime    time.Time
}

func NewSSHModel(userID, username string, width, height int) SSHModel {
	// Load or create progress for this user (by SSH key)
	tracker := engine.NewProgressTracker("qwerty")
	// TODO: Load saved progress from server storage keyed by userID
	
	return SSHModel{
		userID:          userID,
		username:        username,
		width:           width,
		height:          height,
		view:            0,
		typingState:     engine.NewTypingState(""),
		progressTracker: tracker,
		wordGenerator:   engine.NewAdaptiveWordGenerator(tracker),
		keyboard:        components.NewKeyboard(),
	}
}

func (m SSHModel) Init() tea.Cmd {
	return nil
}

func (m SSHModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			if m.view == 1 {
				m.view = 0
			}
			return m, nil
		}

		if m.view == 0 {
			// Welcome view
			if msg.String() == "enter" {
				m.view = 1
				m.startNewSession()
			}
		} else {
			// Practice view
			if len(msg.Runes) > 0 {
				char := msg.Runes[0]
				now := time.Now()
				
				var timeMs int64 = 200
				if !m.lastCharTime.IsZero() {
					timeMs = now.Sub(m.lastCharTime).Milliseconds()
				}
				m.lastCharTime = now

				expectedChar := m.typingState.GetCurrentChar()
				m.typingState.TypeChar(char)
				m.keyboard.SetActiveKey(char)
				
				isCorrect := char == expectedChar
				m.progressTracker.RecordKeypress(expectedChar, isCorrect, timeMs)

				if m.typingState.IsComplete {
					// Auto-continue
					m.startNewSession()
				}
			}
		}
	}

	return m, nil
}

func (m *SSHModel) startNewSession() {
	text := m.wordGenerator.GenerateText(8)
	m.typingState.SetNewText(text)
	m.keyboard.SetActiveKey(0)
	m.lastCharTime = time.Time{}
}

func (m SSHModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.view == 0 {
		return m.renderWelcome()
	}
	return m.renderPractice()
}

func (m SSHModel) renderWelcome() string {
	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Welcome to Hyperkey!")

	userInfo := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render(fmt.Sprintf("Session: %s", m.username))

	progress := m.progressTracker.GetProgressSummary()
	stats := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render(fmt.Sprintf("Letters unlocked: %d/%d", progress.UnlockedCount, progress.TotalLetters))

	prompt := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Bold(true).
		Render("Press Enter to start typing")

	help := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render("q to quit • Your progress is saved by SSH key")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		userInfo,
		stats,
		"",
		prompt,
		"",
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m SSHModel) renderPractice() string {
	progress := m.progressTracker.GetProgressSummary()
	progressInfo := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render(fmt.Sprintf("Letters: %d/%d • Focus: %c",
			progress.UnlockedCount, progress.TotalLetters, progress.FocusLetter))

	// Typing area
	typingArea := components.NewTypingArea(min(m.width-10, 70))
	typingContent := typingArea.Render(m.typingState)

	// Keyboard
	keyboardContent := m.keyboard.Render()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		progressInfo,
		"",
		typingContent,
		"",
		keyboardContent,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
