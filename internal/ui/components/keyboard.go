package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/farseen/hyperkey/internal/ui/styles"
)

// QWERTY keyboard layout
var qwertyRows = [][]string{
	{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
	{"a", "s", "d", "f", "g", "h", "j", "k", "l"},
	{"z", "x", "c", "v", "b", "n", "m"},
}

// Keyboard represents a virtual keyboard visualization
type Keyboard struct {
	activeKey rune
	weakKeys  map[rune]float64
	showHands bool
}

// NewKeyboard creates a new keyboard component
func NewKeyboard() *Keyboard {
	return &Keyboard{
		weakKeys: make(map[rune]float64),
	}
}

// SetActiveKey sets the currently active/pressed key
func (k *Keyboard) SetActiveKey(key rune) {
	k.activeKey = key
}

// SetWeakKeys sets the weak keys map for highlighting
func (k *Keyboard) SetWeakKeys(keys map[rune]float64) {
	k.weakKeys = keys
}

// Simple key style without borders (for horizontal layout)
var (
	simpleKeyStyle = lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Background(styles.BgCard).
		Padding(0, 1).
		MarginRight(1)

	simpleKeyActiveStyle = lipgloss.NewStyle().
		Foreground(styles.BgDark).
		Background(styles.Secondary).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	simpleKeyWeakStyle = lipgloss.NewStyle().
		Foreground(styles.Warning).
		Background(styles.BgCard).
		Padding(0, 1).
		MarginRight(1)
)

// Render renders the keyboard visualization
func (k *Keyboard) Render() string {
	var rows []string

	// Render each row
	for rowIdx, row := range qwertyRows {
		var keys []string

		for _, keyStr := range row {
			key := rune(keyStr[0])
			keyDisplay := strings.ToUpper(keyStr)

			// Choose style based on key state
			var style lipgloss.Style
			if key == k.activeKey || k.activeKey == rune(strings.ToUpper(keyStr)[0]) {
				style = simpleKeyActiveStyle
			} else if _, isWeak := k.weakKeys[key]; isWeak {
				style = simpleKeyWeakStyle
			} else {
				style = simpleKeyStyle
			}

			keys = append(keys, style.Render(keyDisplay))
		}

		// Join keys horizontally
		rowContent := lipgloss.JoinHorizontal(lipgloss.Center, keys...)

		// Add indentation for keyboard stagger
		indent := strings.Repeat(" ", rowIdx*2)
		rows = append(rows, indent+rowContent)
	}

	// Add spacebar
	spaceStyle := simpleKeyStyle.Width(24).Align(lipgloss.Center)
	if k.activeKey == ' ' {
		spaceStyle = simpleKeyActiveStyle.Width(24).Align(lipgloss.Center)
	}
	spaceRow := "      " + spaceStyle.Render("SPACE")
	rows = append(rows, spaceRow)

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// RenderCompact renders a more compact keyboard visualization (no background)
func (k *Keyboard) RenderCompact() string {
	var rows []string

	for rowIdx, row := range qwertyRows {
		var rowContent strings.Builder
		indent := strings.Repeat(" ", rowIdx)
		rowContent.WriteString(indent)

		for _, keyStr := range row {
			key := rune(keyStr[0])
			keyDisplay := keyStr

			var fg lipgloss.Color
			if key == k.activeKey || k.activeKey == rune(strings.ToUpper(keyStr)[0]) {
				fg = styles.Secondary
				keyDisplay = strings.ToUpper(keyStr)
			} else if _, isWeak := k.weakKeys[key]; isWeak {
				fg = styles.Warning
			} else {
				fg = styles.TextMuted
			}

			rowContent.WriteString(lipgloss.NewStyle().Foreground(fg).Render(keyDisplay))
			rowContent.WriteString(" ")
		}
		rows = append(rows, rowContent.String())
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
