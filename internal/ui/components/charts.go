package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/farseenmanekhan1232/hyperkey/internal/engine"
	"github.com/farseenmanekhan1232/hyperkey/internal/ui/styles"
)

// Charts provides ASCII chart rendering for statistics
type Charts struct{}

// NewCharts creates a new charts component
func NewCharts() *Charts {
	return &Charts{}
}

// RenderKeyFrequencyHistogram renders a horizontal bar chart of key frequencies
func (c *Charts) RenderKeyFrequencyHistogram(keyStats map[rune]*engine.KeyStat, maxWidth int) string {
	if len(keyStats) == 0 {
		return lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No data yet")
	}

	// Get all keys and find max value
	type keyStat struct {
		key    rune
		hits   int
		misses int
	}
	var stats []keyStat
	maxHits := 0

	for k, s := range keyStats {
		if s.TotalTyped > 0 {
			hits := s.TotalTyped - s.ErrorCount
			stats = append(stats, keyStat{k, hits, s.ErrorCount})
			if s.TotalTyped > maxHits {
				maxHits = s.TotalTyped
			}
		}
	}

	if len(stats) == 0 {
		return lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No data yet")
	}

	// Sort alphabetically
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].key < stats[j].key
	})

	barWidth := maxWidth - 10 // Leave room for label
	if barWidth < 20 {
		barWidth = 20
	}

	var rows []string
	hitStyle := lipgloss.NewStyle().Foreground(styles.Success)
	missStyle := lipgloss.NewStyle().Foreground(styles.Error)

	for _, s := range stats {
		// Calculate bar lengths
		total := s.hits + s.misses
		hitLen := int(float64(s.hits) / float64(maxHits) * float64(barWidth))
		missLen := int(float64(s.misses) / float64(maxHits) * float64(barWidth))

		if hitLen < 1 && s.hits > 0 {
			hitLen = 1
		}
		if missLen < 1 && s.misses > 0 {
			missLen = 1
		}

		hitBar := hitStyle.Render(strings.Repeat("█", hitLen))
		missBar := missStyle.Render(strings.Repeat("█", missLen))

		label := fmt.Sprintf("%c ", s.key)
		count := fmt.Sprintf(" %d", total)

		rows = append(rows, label+hitBar+missBar+count)
	}

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Key Frequency")

	legend := lipgloss.NewStyle().Foreground(styles.TextMuted).Render(
		"█ hits  " + missStyle.Render("█") + " misses",
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		strings.Join(rows, "\n"),
		"",
		legend,
	)
}

// RenderKeyboardHeatmap renders a keyboard with color-coded keys based on error rate
func (c *Charts) RenderKeyboardHeatmap(keyStats map[rune]*engine.KeyStat) string {
	rows := [][]string{
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l"},
		{"z", "x", "c", "v", "b", "n", "m"},
	}

	// Calculate error rate for each key
	getColor := func(key rune) lipgloss.Color {
		stat, ok := keyStats[key]
		if !ok || stat.TotalTyped < 5 {
			return styles.TextMuted // Not enough data
		}

		errorRate := float64(stat.ErrorCount) / float64(stat.TotalTyped)

		// Color gradient: green (good) -> yellow -> red (bad)
		if errorRate < 0.05 {
			return styles.Success // < 5% errors - excellent
		} else if errorRate < 0.10 {
			return lipgloss.Color("#84CC16") // Light green
		} else if errorRate < 0.15 {
			return styles.Warning // Yellow - needs work
		} else if errorRate < 0.25 {
			return lipgloss.Color("#F97316") // Orange
		}
		return styles.Error // > 25% errors - problem key
	}

	var output []string

	for rowIdx, row := range rows {
		indent := strings.Repeat(" ", rowIdx*2)
		var keys []string

		for _, keyStr := range row {
			key := rune(keyStr[0])
			color := getColor(key)

			keyStyle := lipgloss.NewStyle().
				Foreground(color).
				Background(styles.BgCard).
				Padding(0, 1)

			keys = append(keys, keyStyle.Render(strings.ToUpper(keyStr)))
		}

		output = append(output, indent+lipgloss.JoinHorizontal(lipgloss.Center, keys...))
	}

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Keyboard Heatmap")

	legend := lipgloss.NewStyle().Foreground(styles.TextMuted).Render(
		lipgloss.NewStyle().Foreground(styles.Success).Render("●") + " good  " +
			lipgloss.NewStyle().Foreground(styles.Warning).Render("●") + " ok  " +
			lipgloss.NewStyle().Foreground(styles.Error).Render("●") + " weak",
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		strings.Join(output, "\n"),
		"",
		legend,
	)
}

// RenderSpeedHistogram renders a bar chart of typing speed per key
func (c *Charts) RenderSpeedHistogram(keyStats map[rune]*engine.KeyStat, maxWidth int) string {
	if len(keyStats) == 0 {
		return lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No data yet")
	}

	type keySpeed struct {
		key   rune
		avgMs float64
	}
	var speeds []keySpeed
	maxSpeed := 0.0

	for k, s := range keyStats {
		if s.TotalTyped >= 5 && s.AverageTimeMs > 0 {
			speed := s.AverageTimeMs
			speeds = append(speeds, keySpeed{k, speed})
			if speed > maxSpeed {
				maxSpeed = speed
			}
		}
	}

	if len(speeds) == 0 {
		return lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No speed data yet")
	}

	// Sort alphabetically
	sort.Slice(speeds, func(i, j int) bool {
		return speeds[i].key < speeds[j].key
	})

	barWidth := maxWidth - 15
	if barWidth < 15 {
		barWidth = 15
	}

	var rows []string
	barChar := "▇"

	for _, s := range speeds {
		// Invert: faster is better, so shorter bar for faster times
		// Max bar = slowest, min bar = fastest
		barLen := int((s.avgMs / maxSpeed) * float64(barWidth))
		if barLen < 1 {
			barLen = 1
		}

		// Color based on speed (faster = greener)
		var color lipgloss.Color
		if s.avgMs < 150 {
			color = styles.Success // Fast
		} else if s.avgMs < 250 {
			color = lipgloss.Color("#84CC16")
		} else if s.avgMs < 400 {
			color = styles.Warning
		} else {
			color = styles.Error // Slow
		}

		bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat(barChar, barLen))
		label := fmt.Sprintf("%c ", s.key)
		speed := fmt.Sprintf(" %dms", int(s.avgMs))

		rows = append(rows, label+bar+speed)
	}

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Key Speed (avg time)")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		strings.Join(rows, "\n"),
	)
}

// RenderProgressSummary renders a summary of learning progress
func (c *Charts) RenderProgressSummary(tracker *engine.ProgressTracker) string {
	if tracker == nil {
		return ""
	}

	summary := tracker.GetProgressSummary()

	// Progress bar for letters unlocked
	totalLetters := summary.TotalLetters
	unlocked := summary.UnlockedCount
	progress := float64(unlocked) / float64(totalLetters)
	
	barWidth := 30
	filledLen := int(progress * float64(barWidth))
	emptyLen := barWidth - filledLen

	progressBar := lipgloss.NewStyle().Foreground(styles.Success).Render(strings.Repeat("█", filledLen)) +
		lipgloss.NewStyle().Foreground(styles.TextMuted).Render(strings.Repeat("░", emptyLen))

	progressText := fmt.Sprintf("%d/%d letters (%d%%)", unlocked, totalLetters, int(progress*100))

	// Show unlocked letters
	var unlockedLetters []string
	for _, l := range summary.Unlocked {
		unlockedLetters = append(unlockedLetters, lipgloss.NewStyle().Foreground(styles.Success).Render(string(l)))
	}

	// Show next letters to unlock
	var lockedLetters []string
	layout := engine.GetLayout(summary.LayoutName)
	for i := unlocked; i < min(unlocked+5, len(layout.UnlockOrder)); i++ {
		lockedLetters = append(lockedLetters, lipgloss.NewStyle().Foreground(styles.TextMuted).Render(string(layout.UnlockOrder[i])))
	}

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Learning Progress")

	focus := fmt.Sprintf("Focus: %s", lipgloss.NewStyle().Foreground(styles.Accent).Bold(true).Render(string(summary.FocusLetter)))
	
	unlockRow := "Unlocked: " + strings.Join(unlockedLetters, " ")
	nextRow := "Next: " + strings.Join(lockedLetters, " ")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		progressBar + " " + progressText,
		"",
		focus,
		unlockRow,
		nextRow,
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
