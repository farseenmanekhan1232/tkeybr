package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/farseenmanekhan1232/tkeybr/internal/ui/styles"
)

// Header renders the app header
type Header struct {
	userName  string
	userEmail string
	isGuest   bool
}

// NewHeader creates a new header component
func NewHeader() *Header {
	return &Header{
		isGuest: true,
	}
}

// SetUser sets the logged-in user info
func (h *Header) SetUser(name, email string) {
	h.userName = name
	h.userEmail = email
	h.isGuest = false
}

// SetGuest sets the header to guest mode
func (h *Header) SetGuest() {
	h.userName = ""
	h.userEmail = ""
	h.isGuest = true
}

// Render renders the header
func (h *Header) Render(width int) string {
	// Logo/title
	logo := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("⌨  TKEYBR")

	subtitle := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render(" - Master your keyboard")

	title := logo + subtitle

	// User info
	var userInfo string
	if h.isGuest {
		userInfo = lipgloss.NewStyle().
			Foreground(styles.TextMuted).
			Italic(true).
			Render("Guest mode • Press 'g' to sign in")
	} else {
		userInfo = lipgloss.NewStyle().
			Foreground(styles.TextSecondary).
			Render("👤 " + h.userName)
	}

	// Build header bar
	titleWidth := lipgloss.Width(title)
	userWidth := lipgloss.Width(userInfo)
	spacerWidth := width - titleWidth - userWidth - 4
	if spacerWidth < 0 {
		spacerWidth = 2
	}
	spacer := lipgloss.NewStyle().Width(spacerWidth).Render("")

	content := lipgloss.JoinHorizontal(lipgloss.Center, title, spacer, userInfo)

	return lipgloss.NewStyle().
		Background(styles.BgCard).
		Foreground(styles.TextPrimary).
		Padding(0, 2).
		Width(width).
		Render(content)
}

// RenderMinimal renders a minimal header
func (h *Header) RenderMinimal() string {
	return lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("⌨ TKEYBR")
}
