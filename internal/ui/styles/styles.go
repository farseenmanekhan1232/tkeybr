package styles

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme
type Theme struct {
	Name string

	// Text colors
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color

	// Status colors
	Success lipgloss.Color
	Error   lipgloss.Color
	Warning lipgloss.Color
	Accent  lipgloss.Color

	// Background colors
	BgDark  lipgloss.Color
	BgCard  lipgloss.Color
	BgHover lipgloss.Color

	// Border color
	Border lipgloss.Color
}

// Available themes
var (
	// Minimal - Clean grayscale (default)
	ThemeMinimal = Theme{
		Name:          "minimal",
		TextPrimary:   lipgloss.Color("#F8FAFC"),
		TextSecondary: lipgloss.Color("#94A3B8"),
		TextMuted:     lipgloss.Color("#64748B"),
		Success:       lipgloss.Color("#22C55E"),
		Error:         lipgloss.Color("#EF4444"),
		Warning:       lipgloss.Color("#F59E0B"),
		Accent:        lipgloss.Color("#F59E0B"),
		BgDark:        lipgloss.Color("#0F172A"),
		BgCard:        lipgloss.Color("#1E293B"),
		BgHover:       lipgloss.Color("#334155"),
		Border:        lipgloss.Color("#475569"),
	}

	// Dracula - Popular dark theme
	ThemeDracula = Theme{
		Name:          "dracula",
		TextPrimary:   lipgloss.Color("#F8F8F2"),
		TextSecondary: lipgloss.Color("#6272A4"),
		TextMuted:     lipgloss.Color("#44475A"),
		Success:       lipgloss.Color("#50FA7B"),
		Error:         lipgloss.Color("#FF5555"),
		Warning:       lipgloss.Color("#FFB86C"),
		Accent:        lipgloss.Color("#BD93F9"),
		BgDark:        lipgloss.Color("#282A36"),
		BgCard:        lipgloss.Color("#44475A"),
		BgHover:       lipgloss.Color("#6272A4"),
		Border:        lipgloss.Color("#6272A4"),
	}

	// Nord - Cool blue tones
	ThemeNord = Theme{
		Name:          "nord",
		TextPrimary:   lipgloss.Color("#ECEFF4"),
		TextSecondary: lipgloss.Color("#D8DEE9"),
		TextMuted:     lipgloss.Color("#4C566A"),
		Success:       lipgloss.Color("#A3BE8C"),
		Error:         lipgloss.Color("#BF616A"),
		Warning:       lipgloss.Color("#EBCB8B"),
		Accent:        lipgloss.Color("#88C0D0"),
		BgDark:        lipgloss.Color("#2E3440"),
		BgCard:        lipgloss.Color("#3B4252"),
		BgHover:       lipgloss.Color("#434C5E"),
		Border:        lipgloss.Color("#4C566A"),
	}

	// Monokai - Classic coding theme
	ThemeMonokai = Theme{
		Name:          "monokai",
		TextPrimary:   lipgloss.Color("#F8F8F2"),
		TextSecondary: lipgloss.Color("#75715E"),
		TextMuted:     lipgloss.Color("#49483E"),
		Success:       lipgloss.Color("#A6E22E"),
		Error:         lipgloss.Color("#F92672"),
		Warning:       lipgloss.Color("#E6DB74"),
		Accent:        lipgloss.Color("#66D9EF"),
		BgDark:        lipgloss.Color("#272822"),
		BgCard:        lipgloss.Color("#3E3D32"),
		BgHover:       lipgloss.Color("#49483E"),
		Border:        lipgloss.Color("#75715E"),
	}

	// Solarized Dark
	ThemeSolarized = Theme{
		Name:          "solarized",
		TextPrimary:   lipgloss.Color("#839496"),
		TextSecondary: lipgloss.Color("#657B83"),
		TextMuted:     lipgloss.Color("#586E75"),
		Success:       lipgloss.Color("#859900"),
		Error:         lipgloss.Color("#DC322F"),
		Warning:       lipgloss.Color("#B58900"),
		Accent:        lipgloss.Color("#268BD2"),
		BgDark:        lipgloss.Color("#002B36"),
		BgCard:        lipgloss.Color("#073642"),
		BgHover:       lipgloss.Color("#094554"),
		Border:        lipgloss.Color("#586E75"),
	}

	// Gruvbox - Retro warm theme  
	ThemeGruvbox = Theme{
		Name:          "gruvbox",
		TextPrimary:   lipgloss.Color("#EBDBB2"),
		TextSecondary: lipgloss.Color("#A89984"),
		TextMuted:     lipgloss.Color("#665C54"),
		Success:       lipgloss.Color("#B8BB26"),
		Error:         lipgloss.Color("#FB4934"),
		Warning:       lipgloss.Color("#FABD2F"),
		Accent:        lipgloss.Color("#83A598"),
		BgDark:        lipgloss.Color("#282828"),
		BgCard:        lipgloss.Color("#3C3836"),
		BgHover:       lipgloss.Color("#504945"),
		Border:        lipgloss.Color("#665C54"),
	}

	// All available themes
	Themes = map[string]Theme{
		"minimal":   ThemeMinimal,
		"dracula":   ThemeDracula,
		"nord":      ThemeNord,
		"monokai":   ThemeMonokai,
		"solarized": ThemeSolarized,
		"gruvbox":   ThemeGruvbox,
	}

	// ThemeNames for display
	ThemeNames = []string{"minimal", "dracula", "nord", "monokai", "solarized", "gruvbox"}
)

// Current active theme
var CurrentTheme = ThemeMinimal

// SetTheme changes the active theme
func SetTheme(name string) {
	if theme, ok := Themes[name]; ok {
		CurrentTheme = theme
		updateStyles()
	}
}

// GetTheme returns the current theme
func GetTheme() Theme {
	return CurrentTheme
}

// updateStyles rebuilds all styles with current theme
func updateStyles() {
	t := CurrentTheme

	HeaderStyle = lipgloss.NewStyle().
		Foreground(t.TextPrimary).
		Background(t.BgCard).
		Padding(0, 2).
		Bold(true)

	TitleStyle = lipgloss.NewStyle().
		Foreground(t.TextPrimary).
		Bold(true).
		Padding(0, 1)

	StatsContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(1, 2)

	StatLabelStyle = lipgloss.NewStyle().
		Foreground(t.TextSecondary)

	StatValueStyle = lipgloss.NewStyle().
		Foreground(t.TextPrimary).
		Bold(true)

	TypingAreaStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(1, 2).
		Width(80)

	CorrectCharStyle = lipgloss.NewStyle().
		Foreground(t.Success)

	IncorrectCharStyle = lipgloss.NewStyle().
		Foreground(t.Error).
		Underline(true)

	CurrentCharStyle = lipgloss.NewStyle().
		Foreground(t.BgDark).
		Background(t.Accent)

	PendingCharStyle = lipgloss.NewStyle().
		Foreground(t.TextMuted)

	KeyStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(0, 1).
		Foreground(t.TextSecondary)

	KeyActiveStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Background(t.Accent).
		Padding(0, 1).
		Foreground(t.BgDark).
		Bold(true)

	KeyWeakStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Warning).
		Padding(0, 1).
		Foreground(t.Warning)

	HelpStyle = lipgloss.NewStyle().
		Foreground(t.TextMuted).
		Padding(1, 0)

	AdSlotStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Border).
		Foreground(t.TextMuted).
		Padding(0, 2).
		Align(lipgloss.Center)

	MenuItemStyle = lipgloss.NewStyle().
		Foreground(t.TextSecondary).
		Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
		Foreground(t.TextPrimary).
		Background(t.BgHover).
		Padding(0, 2).
		Bold(true)

	UserInfoStyle = lipgloss.NewStyle().
		Foreground(t.TextSecondary).
		Italic(true)

	AttentionStyle = lipgloss.NewStyle().
		Foreground(t.BgDark).
		Background(t.Accent).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center)
}

// Expose colors from current theme
var (
	TextPrimary   = ThemeMinimal.TextPrimary
	TextSecondary = ThemeMinimal.TextSecondary
	TextMuted     = ThemeMinimal.TextMuted
	Success       = ThemeMinimal.Success
	Error         = ThemeMinimal.Error
	Warning       = ThemeMinimal.Warning
	Accent        = ThemeMinimal.Accent
	BgDark        = ThemeMinimal.BgDark
	BgCard        = ThemeMinimal.BgCard
	BgHover       = ThemeMinimal.BgHover
	Border        = ThemeMinimal.Border
	Primary       = ThemeMinimal.TextSecondary
)

// Component styles (initialized with default theme)
var (
	HeaderStyle           lipgloss.Style
	TitleStyle            lipgloss.Style
	StatsContainerStyle   lipgloss.Style
	StatLabelStyle        lipgloss.Style
	StatValueStyle        lipgloss.Style
	TypingAreaStyle       lipgloss.Style
	CorrectCharStyle      lipgloss.Style
	IncorrectCharStyle    lipgloss.Style
	CurrentCharStyle      lipgloss.Style
	PendingCharStyle      lipgloss.Style
	KeyStyle              lipgloss.Style
	KeyActiveStyle        lipgloss.Style
	KeyWeakStyle          lipgloss.Style
	HelpStyle             lipgloss.Style
	AdSlotStyle           lipgloss.Style
	MenuItemStyle         lipgloss.Style
	MenuItemSelectedStyle lipgloss.Style
	UserInfoStyle         lipgloss.Style
	AttentionStyle        lipgloss.Style
)

func init() {
	updateStyles()
}

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
