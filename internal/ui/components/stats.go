package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/farseenmanekhan1232/hyperkey/internal/engine"
	"github.com/farseenmanekhan1232/hyperkey/internal/ui/styles"
)

// StatsPanel renders typing statistics
type StatsPanel struct {
	compact bool
}

// NewStatsPanel creates a new stats panel component
func NewStatsPanel(compact bool) *StatsPanel {
	return &StatsPanel{compact: compact}
}

// Render renders the stats panel
func (sp *StatsPanel) Render(stats *engine.Stats) string {
	if stats == nil {
		return ""
	}

	if sp.compact {
		return sp.renderCompact(stats)
	}
	return sp.renderFull(stats)
}

func (sp *StatsPanel) renderCompact(stats *engine.Stats) string {
	wpmLabel := styles.StatLabelStyle.Render("WPM")
	wpmValue := styles.StatValueStyle.Render(fmt.Sprintf("%.0f", stats.CurrentWPM))
	
	accLabel := styles.StatLabelStyle.Render("ACC")
	accValue := styles.StatValueStyle.Render(fmt.Sprintf("%.0f%%", stats.CurrentAccuracy))

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("%s %s", wpmLabel, wpmValue),
		"  │  ",
		fmt.Sprintf("%s %s", accLabel, accValue),
	)
}

func (sp *StatsPanel) renderFull(stats *engine.Stats) string {
	// WPM box
	wpm := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.StatValueStyle.Render(fmt.Sprintf("%.0f", stats.CurrentWPM)),
		styles.StatLabelStyle.Render("WPM"),
	)
	wpmBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(0, 2).
		Render(wpm)

	// Accuracy box
	accuracy := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.StatValueStyle.Render(fmt.Sprintf("%.0f%%", stats.CurrentAccuracy)),
		styles.StatLabelStyle.Render("Accuracy"),
	)
	accuracyBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Success).
		Padding(0, 2).
		Render(accuracy)

	// Characters box
	chars := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.StatValueStyle.Render(fmt.Sprintf("%d", stats.CharactersTyped)),
		styles.StatLabelStyle.Render("Chars"),
	)
	charsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Secondary).
		Padding(0, 2).
		Render(chars)

	// Errors box
	errColor := styles.TextMuted
	if stats.ErrorCount > 0 {
		errColor = styles.Error
	}
	errors := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(errColor).Bold(true).Render(fmt.Sprintf("%d", stats.ErrorCount)),
		styles.StatLabelStyle.Render("Errors"),
	)
	errorsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errColor).
		Padding(0, 2).
		Render(errors)

	return lipgloss.JoinHorizontal(lipgloss.Center, wpmBox, " ", accuracyBox, " ", charsBox, " ", errorsBox)
}

// RenderHistory renders historical stats
func (sp *StatsPanel) RenderHistory(stats *engine.Stats) string {
	if stats == nil || stats.TotalSessions == 0 {
		return styles.StatLabelStyle.Render("No history yet")
	}

	lines := []string{
		fmt.Sprintf("Sessions: %d", stats.TotalSessions),
		fmt.Sprintf("Best WPM: %.0f", stats.BestWPM),
		fmt.Sprintf("Avg WPM: %.0f", stats.AverageWPM),
		fmt.Sprintf("Total chars: %d", stats.TotalCharacters),
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.TextMuted).
		Padding(0, 1).
		Render(content)
}
