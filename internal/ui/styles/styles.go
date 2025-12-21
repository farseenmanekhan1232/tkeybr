package styles

import "github.com/charmbracelet/lipgloss"

// Color palette - modern dark theme
var (
	// Primary colors
	Primary     = lipgloss.Color("#7C3AED") // Purple
	Secondary   = lipgloss.Color("#06B6D4") // Cyan
	Accent      = lipgloss.Color("#F59E0B") // Amber
	
	// Text colors
	TextPrimary   = lipgloss.Color("#F8FAFC")
	TextSecondary = lipgloss.Color("#94A3B8")
	TextMuted     = lipgloss.Color("#64748B")
	
	// Status colors
	Success = lipgloss.Color("#22C55E") // Green
	Error   = lipgloss.Color("#EF4444") // Red
	Warning = lipgloss.Color("#F59E0B") // Amber
	
	// Background colors
	BgDark    = lipgloss.Color("#0F172A")
	BgCard    = lipgloss.Color("#1E293B")
	BgHover   = lipgloss.Color("#334155")
	BgKeyboard = lipgloss.Color("#1E293B")
)

// Component styles
var (
	// Header style
	HeaderStyle = lipgloss.NewStyle().
		Foreground(TextPrimary).
		Background(BgCard).
		Padding(0, 2).
		Bold(true)

	// Title style
	TitleStyle = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		Padding(0, 1)

	// Stats container
	StatsContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2)

	// Stat label
	StatLabelStyle = lipgloss.NewStyle().
		Foreground(TextSecondary)

	// Stat value
	StatValueStyle = lipgloss.NewStyle().
		Foreground(TextPrimary).
		Bold(true)

	// Typing area
	TypingAreaStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Secondary).
		Padding(1, 2).
		Width(80)

	// Correct character
	CorrectCharStyle = lipgloss.NewStyle().
		Foreground(Success)

	// Incorrect character
	IncorrectCharStyle = lipgloss.NewStyle().
		Foreground(Error).
		Underline(true)

	// Current character (cursor)
	CurrentCharStyle = lipgloss.NewStyle().
		Foreground(BgDark).
		Background(Accent)

	// Pending character (not yet typed)
	PendingCharStyle = lipgloss.NewStyle().
		Foreground(TextMuted)

	// Keyboard key styles
	KeyStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(TextMuted).
		Padding(0, 1).
		Foreground(TextSecondary)

	KeyActiveStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Secondary).
		Background(Secondary).
		Padding(0, 1).
		Foreground(BgDark).
		Bold(true)

	KeyWeakStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Warning).
		Padding(0, 1).
		Foreground(Warning)

	// Footer/help style
	HelpStyle = lipgloss.NewStyle().
		Foreground(TextMuted).
		Padding(1, 0)

	// Ad slot placeholder
	AdSlotStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(TextMuted).
		Foreground(TextMuted).
		Padding(0, 2).
		Align(lipgloss.Center)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
		Foreground(TextSecondary).
		Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
		Foreground(Primary).
		Background(BgHover).
		Padding(0, 2).
		Bold(true)

	// User info style
	UserInfoStyle = lipgloss.NewStyle().
		Foreground(TextSecondary).
		Italic(true)
)

// Layout helpers
func Center(width int, content string) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(content)
}

func Container(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2)
}
