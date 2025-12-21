package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/farseen/hyperkey/internal/engine"
	"github.com/farseen/hyperkey/internal/ui/styles"
)

// TypingArea renders the typing practice area
type TypingArea struct {
	width int
}

// NewTypingArea creates a new typing area component
func NewTypingArea(width int) *TypingArea {
	return &TypingArea{width: width}
}

// Render renders the typing area with the given typing state
func (ta *TypingArea) Render(ts *engine.TypingState) string {
	if ts == nil || ts.TargetText == "" {
		return styles.TypingAreaStyle.Render("Loading...")
	}

	var output strings.Builder

	// Render the text with character-by-character styling
	for i, char := range ts.TargetText {
		status := ts.GetCharStatus(i)
		charStr := string(char)
		
		var style lipgloss.Style
		switch status {
		case "correct":
			style = styles.CorrectCharStyle
		case "incorrect":
			style = styles.IncorrectCharStyle
		case "current":
			style = styles.CurrentCharStyle
		case "pending":
			style = styles.PendingCharStyle
		}

		output.WriteString(style.Render(charStr))
	}

	// Wrap in container
	content := output.String()
	
	// Word wrap for long text
	wrapped := ta.wordWrap(content, ta.width-4)
	
	return styles.TypingAreaStyle.Width(ta.width).Render(wrapped)
}

// wordWrap wraps text at word boundaries
func (ta *TypingArea) wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	// For styled text, we can't easily word wrap
	// Just return as is for now - the container will handle overflow
	return text
}

// RenderMinimal renders a minimal typing area without borders
func (ta *TypingArea) RenderMinimal(ts *engine.TypingState) string {
	if ts == nil || ts.TargetText == "" {
		return "Preparing text..."
	}

	var output strings.Builder
	lineWidth := 0
	maxWidth := ta.width - 4

	for i, char := range ts.TargetText {
		status := ts.GetCharStatus(i)
		charStr := string(char)
		
		var style lipgloss.Style
		switch status {
		case "correct":
			style = styles.CorrectCharStyle
		case "incorrect":
			style = styles.IncorrectCharStyle
		case "current":
			style = styles.CurrentCharStyle
		case "pending":
			style = styles.PendingCharStyle
		}

		// Handle word wrap at spaces
		if char == ' ' && lineWidth >= maxWidth-10 {
			output.WriteString("\n")
			lineWidth = 0
		}

		output.WriteString(style.Render(charStr))
		lineWidth++
	}

	return output.String()
}
