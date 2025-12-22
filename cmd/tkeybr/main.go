package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/farseenmanekhan1232/tkeybr/internal/app"
)

// Build-time variables (set via ldflags)
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("tkeybr %s (%s)\n", version, commit)
		return
	}

	p := tea.NewProgram(
		app.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

