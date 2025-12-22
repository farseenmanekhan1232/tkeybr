package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/farseenmanekhan1232/tkeybr/internal/ui/styles"
)

// HelpScreen renders the help/documentation screen
type HelpScreen struct {
	width     int
	scrollPos int
}

// NewHelpScreen creates a new help screen
func NewHelpScreen() *HelpScreen {
	return &HelpScreen{}
}

// SetWidth sets the width for the help screen
func (h *HelpScreen) SetWidth(width int) {
	h.width = width
}

// ScrollUp scrolls the help content up
func (h *HelpScreen) ScrollUp() {
	if h.scrollPos > 0 {
		h.scrollPos--
	}
}

// ScrollDown scrolls the help content down
func (h *HelpScreen) ScrollDown() {
	h.scrollPos++
}

// Render renders the help screen content
func (h *HelpScreen) Render() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Bold(true).
		MarginTop(1)

	textStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary)

	accentStyle := lipgloss.NewStyle().
		Foreground(styles.Accent)

	successStyle := lipgloss.NewStyle().
		Foreground(styles.Success)

	// Build help content
	content := []string{
		titleStyle.Render("📖 How tkeybr Works"),
		"",
		sectionStyle.Render("Touch Typing"),
		textStyle.Render("tkeybr helps you learn touch typing - typing through"),
		textStyle.Render("muscle memory without looking at the keyboard."),
		"",
		sectionStyle.Render("The Teaching Method"),
		textStyle.Render("We use a unique adaptive approach with smart algorithms:"),
		"",
		accentStyle.Render("1. Generate") + textStyle.Render(" → Words based on your skill level"),
		accentStyle.Render("2. Practice") + textStyle.Render(" → Type with minimal errors"),
		accentStyle.Render("3. Analyze") + textStyle.Render("  → Track your per-key statistics"),
		accentStyle.Render("4. Adapt") + textStyle.Render("    → Generate better words for you"),
		"",
		sectionStyle.Render("Progressive Letter Unlocking"),
		textStyle.Render("You start with a small set of letters (home row):"),
		successStyle.Render("  f j d k") + textStyle.Render(" ← Your starting letters"),
		"",
		textStyle.Render("As you improve, more letters unlock:"),
		successStyle.Render("  f j d k s l g h...") + textStyle.Render(" and so on"),
		"",
		sectionStyle.Render("The Focus Letter"),
		textStyle.Render("One letter is always your ") + accentStyle.Render("focus letter") + textStyle.Render(" - the key"),
		textStyle.Render("you struggle with most. It appears in every word!"),
		"",
		sectionStyle.Render("Keyboard Zones"),
		textStyle.Render("Each finger is responsible for specific keys:"),
		"",
		lipgloss.NewStyle().Foreground(styles.TextMuted).Render("  Left Hand        Right Hand"),
		lipgloss.NewStyle().Foreground(styles.TextMuted).Render("  ──────────       ──────────"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("  Q A Z") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316")).Render("W S X") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308")).Render("E D C") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Render("R T F G V B") + "   " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Render("Y U H J N M") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308")).Render("I K") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316")).Render("O L") + "  " +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("P"),
		lipgloss.NewStyle().Foreground(styles.TextMuted).Render("  Pinky Ring Mid  Index        Index        Mid Ring Pinky"),
		"",
		sectionStyle.Render("Tips for Success"),
		"",
		successStyle.Render("  ✓") + textStyle.Render(" Keep fingers on home row (F and J bumps)"),
		successStyle.Render("  ✓") + textStyle.Render(" Don't look at the keyboard"),
		successStyle.Render("  ✓") + textStyle.Render(" Accuracy over speed - errors hurt more"),
		successStyle.Render("  ✓") + textStyle.Render(" Practice daily, even 10 minutes helps"),
		successStyle.Render("  ✓") + textStyle.Render(" Be patient - progress takes time"),
		"",
		sectionStyle.Render("Expected Progress"),
		textStyle.Render("Real user examples:"),
		textStyle.Render("  • 30→70 WPM in ~4 hours over 15 days"),
		textStyle.Render("  • 35→70 WPM in ~2.5 hours over 12 days"),
		textStyle.Render("  • 20→45 WPM in ~10 hours over 22 days"),
		"",
		lipgloss.NewStyle().Foreground(styles.TextMuted).Italic(true).Render("Press Esc to go back"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}
