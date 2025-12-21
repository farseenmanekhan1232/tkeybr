package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/farseen/hyperkey/internal/auth"
	"github.com/farseen/hyperkey/internal/engine"
	"github.com/farseen/hyperkey/internal/storage"
	"github.com/farseen/hyperkey/internal/ui/components"
	"github.com/farseen/hyperkey/internal/ui/styles"
)

// View represents different app views
type View int

const (
	ViewMenu View = iota
	ViewPractice
	ViewStats
	ViewSettings
	ViewAuth
	ViewHelp
)

// Model represents the application state
type Model struct {
	// View state
	currentView View
	width       int
	height      int

	// Typing state
	typingState     *engine.TypingState
	progressTracker *engine.ProgressTracker
	wordGenerator   *engine.AdaptiveWordGenerator
	stats           *engine.Stats
	lastCharTime    time.Time

	// Components
	header     *components.Header
	keyboard   *components.Keyboard
	typingArea *components.TypingArea
	statsPanel *components.StatsPanel

	// Auth
	auth           *auth.GoogleAuth
	deviceCodeResp *auth.DeviceCodeResponse
	authPolling    bool
	authError      string

	// Storage
	storage     *storage.LocalStorage
	preferences *storage.UserPreferences

	// Menu state
	menuSelection int
	menuItems     []string

	// Ticker for WPM updates
	tickerActive bool
}

// New creates a new application model
func New() Model {
	// Initialize storage and load data
	store := storage.NewLocalStorage()
	savedStats, _ := store.LoadStats()
	prefs, _ := store.LoadPreferences()

	if savedStats == nil {
		savedStats = engine.NewStats()
	}
	if prefs == nil {
		prefs = storage.DefaultPreferences()
	}

	// Initialize progress tracker with layout
	progressTracker := engine.NewProgressTracker(prefs.KeyboardLayout)
	
	// Initialize adaptive word generator
	wordGen := engine.NewAdaptiveWordGenerator(progressTracker)

	// Initialize auth (placeholder client ID - user needs to configure)
	googleAuth := auth.NewGoogleAuth("YOUR_GOOGLE_CLIENT_ID")
	
	// Try to load existing auth
	googleAuth.LoadAuth()

	// Initialize components
	header := components.NewHeader()
	if user := googleAuth.GetCurrentUser(); user != nil {
		header.SetUser(user.Name, user.Email)
	}

	return Model{
		currentView:     ViewMenu,
		typingState:     engine.NewTypingState(""),
		progressTracker: progressTracker,
		wordGenerator:   wordGen,
		stats:           savedStats,
		header:          header,
		keyboard:        components.NewKeyboard(),
		typingArea:      components.NewTypingArea(70),
		statsPanel:      components.NewStatsPanel(false),
		auth:            googleAuth,
		storage:         store,
		preferences:     prefs,
		menuSelection:   0,
		menuItems:       []string{"Start Practice", "View Statistics", "Help", "Settings", "Sign In with Google", "Quit"},
	}
}

// Messages
type tickMsg time.Time
type authPollMsg struct{}
type authCompleteMsg struct {
	user *auth.UserInfo
	err  error
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

// Update handles messages and input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.typingArea = components.NewTypingArea(min(70, msg.Width-10))
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tickMsg:
		if m.currentView == ViewPractice && !m.typingState.IsComplete {
			m.stats.UpdateFromTypingState(m.typingState)
			return m, tickCmd()
		}
		return m, nil

	case authPollMsg:
		return m.handleAuthPoll()

	case authCompleteMsg:
		m.authPolling = false
		if msg.err != nil {
			m.authError = msg.err.Error()
		} else if msg.user != nil {
			m.header.SetUser(msg.user.Name, msg.user.Email)
			m.currentView = ViewMenu
			// Update menu items
			m.menuItems[3] = "Sign Out"
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch {
	case key.Matches(msg, keys.Quit):
		m.saveProgress()
		return m, tea.Quit
	case key.Matches(msg, keys.Escape):
		if m.currentView != ViewMenu {
			m.saveProgress()
			m.currentView = ViewMenu
			m.tickerActive = false
			return m, nil
		}
	}

	switch m.currentView {
	case ViewMenu:
		return m.handleMenuInput(msg)
	case ViewPractice:
		return m.handlePracticeInput(msg)
	case ViewAuth:
		return m.handleAuthInput(msg)
	case ViewStats:
		return m.handleStatsInput(msg)
	case ViewSettings:
		return m.handleSettingsInput(msg)
	}

	return m, nil
}

func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuSelection > 0 {
			m.menuSelection--
		}
	case "down", "j":
		if m.menuSelection < len(m.menuItems)-1 {
			m.menuSelection++
		}
	case "enter":
		switch m.menuSelection {
		case 0: // Start Practice
			m.startNewSession()
			m.currentView = ViewPractice
			m.tickerActive = true
			return m, tickCmd()
		case 1: // View Statistics
			m.currentView = ViewStats
		case 2: // Help
			m.currentView = ViewHelp
		case 3: // Settings
			m.currentView = ViewSettings
		case 4: // Sign In / Sign Out
			if m.auth.GetCurrentUser() != nil {
				m.auth.Logout()
				m.header.SetGuest()
				m.menuItems[4] = "Sign In with Google"
			} else {
				m.currentView = ViewAuth
				return m, m.startAuthFlow()
			}
		case 5: // Quit
			m.saveProgress()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) handlePracticeInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle special keys
	switch msg.Type {
	case tea.KeyBackspace:
		// No backspace allowed - face your mistakes! (keybr style)
		return m, nil
	case tea.KeyTab, tea.KeyEnter:
		// Restart on tab or complete
		if m.typingState.IsComplete {
			m.stats.CompleteSession(m.typingState)
			m.startNewSession()
			return m, tickCmd()
		}
		return m, nil
	}

	// Handle character input
	if len(msg.Runes) > 0 {
		char := msg.Runes[0]
		now := time.Now()
		
		// Calculate time since last keypress
		var timeMs int64 = 200 // Default for first char
		if !m.lastCharTime.IsZero() {
			timeMs = now.Sub(m.lastCharTime).Milliseconds()
		}
		m.lastCharTime = now

		// Get current expected char before typing
		expectedChar := m.typingState.GetCurrentChar()
		
		// Record the keypress
		m.typingState.TypeChar(char)
		m.keyboard.SetActiveKey(char)
		m.stats.UpdateFromTypingState(m.typingState)
		
		// Update progress tracker with timing info
		isCorrect := char == expectedChar
		m.progressTracker.RecordKeypress(expectedChar, isCorrect, timeMs)
		
		// Update weak keys display
		weakKeys := m.progressTracker.GetLetterWeights()
		weakKeyMap := make(map[rune]float64)
		for k, v := range weakKeys {
			if v > 1.3 { // Only show notably weak keys
				weakKeyMap[k] = v
			}
		}
		m.keyboard.SetWeakKeys(weakKeyMap)

		if m.typingState.IsComplete {
			m.stats.CompleteSession(m.typingState)
			m.saveProgress()
			// Auto-start new session - no need to press Tab!
			m.startNewSession()
		}
	}

	return m, nil
}

func (m Model) handleAuthInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Cancel auth on escape (handled globally)
	return m, nil
}

func (m Model) handleStatsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Stats view just needs escape to go back
	return m, nil
}

func (m Model) handleSettingsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Settings navigation would go here
	return m, nil
}

func (m *Model) startNewSession() {
	// Generate new text using adaptive word generator
	text := m.wordGenerator.GenerateText(m.preferences.WordCount)
	
	// Update weak keys on keyboard
	weakKeys := m.progressTracker.GetLetterWeights()
	weakKeyMap := make(map[rune]float64)
	for k, v := range weakKeys {
		if v > 1.3 {
			weakKeyMap[k] = v
		}
	}
	m.keyboard.SetWeakKeys(weakKeyMap)

	m.typingState.SetNewText(text)
	m.keyboard.SetActiveKey(0)
	m.lastCharTime = time.Time{}
}

func (m *Model) saveProgress() {
	m.storage.SaveStats(m.stats)
	m.storage.SavePreferences(m.preferences)
}

func (m *Model) startAuthFlow() tea.Cmd {
	return func() tea.Msg {
		if !m.auth.IsConfigured() {
			return authCompleteMsg{err: fmt.Errorf("Google OAuth not configured. Please set your client ID.")}
		}

		resp, err := m.auth.RequestDeviceCode()
		if err != nil {
			return authCompleteMsg{err: err}
		}

		m.deviceCodeResp = resp
		m.authPolling = true
		return authPollMsg{}
	}
}

func (m Model) handleAuthPoll() (tea.Model, tea.Cmd) {
	if m.deviceCodeResp == nil || !m.authPolling {
		return m, nil
	}

	// Poll for token
	tokenResp, err := m.auth.PollForToken(m.deviceCodeResp.DeviceCode, m.deviceCodeResp.Interval)
	if err != nil {
		return m, tea.Cmd(func() tea.Msg {
			return authCompleteMsg{err: err}
		})
	}

	if tokenResp.Error == "authorization_pending" {
		// Keep polling
		return m, tea.Cmd(func() tea.Msg {
			time.Sleep(time.Duration(m.deviceCodeResp.Interval) * time.Second)
			return authPollMsg{}
		})
	}

	if tokenResp.Error != "" {
		return m, tea.Cmd(func() tea.Msg {
			return authCompleteMsg{err: fmt.Errorf(tokenResp.ErrorDesc)}
		})
	}

	// Success - complete auth
	err = m.auth.CompleteAuth(tokenResp)
	if err != nil {
		return m, tea.Cmd(func() tea.Msg {
			return authCompleteMsg{err: err}
		})
	}

	return m, tea.Cmd(func() tea.Msg {
		return authCompleteMsg{user: m.auth.GetCurrentUser()}
	})
}

// View renders the application
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var content string

	switch m.currentView {
	case ViewMenu:
		content = m.renderMenu()
	case ViewPractice:
		content = m.renderPractice()
	case ViewStats:
		content = m.renderStats()
	case ViewSettings:
		content = m.renderSettings()
	case ViewAuth:
		content = m.renderAuth()
	case ViewHelp:
		content = m.renderHelpView()
	}

	// Build full layout with responsive sizing
	header := m.header.Render(m.width)
	help := m.renderHelp()
	
	// Ad slot at bottom
	adSlot := styles.AdSlotStyle.Width(m.width - 4).Render(
		"[ Ad Space - Support Hyperkey ]",
	)

	// Calculate available height for content
	headerHeight := 1
	helpHeight := 1
	adHeight := 1
	padding := 4
	contentHeight := m.height - headerHeight - helpHeight - adHeight - padding
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Place content in available space
	mainContent := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// Join all sections vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainContent,
		styles.HelpStyle.Width(m.width).Align(lipgloss.Center).Render(help),
		adSlot,
	)
}

func (m Model) renderMenu() string {
	var items []string

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(1).
		Render("Welcome to Hyperkey!")

	subtitle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		MarginBottom(2).
		Render("Master your keyboard with adaptive practice")

	for i, item := range m.menuItems {
		var style lipgloss.Style
		if i == m.menuSelection {
			style = styles.MenuItemSelectedStyle
			item = "▸ " + item
		} else {
			style = styles.MenuItemStyle
			item = "  " + item
		}
		items = append(items, style.Render(item))
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, items...)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		menu,
	)
}

func (m Model) renderPractice() string {
	// Stats at top
	statsRow := m.statsPanel.Render(m.stats)
	
	// Progress info
	progress := m.progressTracker.GetProgressSummary()
	progressInfo := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render(fmt.Sprintf("Letters: %d/%d unlocked • Focus: %c • Layout: %s",
			progress.UnlockedCount, progress.TotalLetters,
			progress.FocusLetter, strings.ToUpper(progress.LayoutName)))

	// Typing area - make width responsive
	typingWidth := min(m.width-10, 80)
	if typingWidth < 40 {
		typingWidth = 40
	}
	typingArea := components.NewTypingArea(typingWidth)
	typingContent := typingArea.Render(m.typingState)

	// Keyboard
	var keyboardContent string
	if m.preferences.ShowKeyboard {
		keyboardContent = m.keyboard.Render()
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		statsRow,
		progressInfo,
		"",
		typingContent,
		"",
		keyboardContent,
	)
}

func (m Model) renderStats() string {
	charts := components.NewCharts()

	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(1).
		Render("📊 Statistics Dashboard")

	// Basic stats
	history := m.statsPanel.RenderHistory(m.stats)

	// Progress summary
	progressSummary := charts.RenderProgressSummary(m.progressTracker)

	// Keyboard heatmap
	heatmap := charts.RenderKeyboardHeatmap(m.stats.KeyStats)

	// Key frequency histogram (compact)
	histogram := charts.RenderKeyFrequencyHistogram(m.stats.KeyStats, 40)

	// Layout: two columns
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		history,
		"",
		progressSummary,
	)

	rightColumn := lipgloss.JoinVertical(lipgloss.Center,
		heatmap,
	)

	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().MarginRight(4).Render(leftColumn),
		rightColumn,
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		columns,
		"",
		histogram,
	)
}

func (m Model) renderSettings() string {
	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(2).
		Render("Settings")

	settings := []string{
		fmt.Sprintf("Keyboard Layout: %s", strings.ToUpper(m.preferences.KeyboardLayout)),
		fmt.Sprintf("Show Keyboard: %v", m.preferences.ShowKeyboard),
		fmt.Sprintf("Word Count: %d", m.preferences.WordCount),
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		strings.Join(settings, "\n"),
		"",
		styles.HelpStyle.Render("Press Esc to go back"),
	)
}

func (m Model) renderAuth() string {
	title := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(2).
		Render("Sign In with Google")

	var content string

	if m.authError != "" {
		content = lipgloss.NewStyle().
			Foreground(styles.Error).
			Render("Error: " + m.authError)
	} else if m.deviceCodeResp != nil {
		instructions := []string{
			"To sign in:",
			"",
			fmt.Sprintf("1. Go to: %s", m.deviceCodeResp.VerificationURL),
			fmt.Sprintf("2. Enter code: %s", lipgloss.NewStyle().Bold(true).Foreground(styles.Accent).Render(m.deviceCodeResp.UserCode)),
			"",
			"Waiting for authorization...",
		}
		content = strings.Join(instructions, "\n")
	} else if !m.auth.IsConfigured() {
		content = lipgloss.NewStyle().
			Foreground(styles.Warning).
			Render("Google OAuth not configured.\n\nPlease add your Google Client ID to the auth configuration.")
	} else {
		content = "Initiating sign-in..."
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		content,
		"",
		styles.HelpStyle.Render("Press Esc to cancel"),
	)
}

func (m Model) renderHelpView() string {
	helpScreen := components.NewHelpScreen()
	helpScreen.SetWidth(m.width)
	return helpScreen.Render()
}

func (m Model) renderHelp() string {
	var help string

	switch m.currentView {
	case ViewMenu:
		help = "↑/↓ Navigate • Enter Select • q Quit"
	case ViewPractice:
		help = "Type to practice • Esc Menu"
	default:
		help = "Esc Back • q Quit"
	}

	return styles.HelpStyle.Render(help)
}

// Tick command for updating WPM
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Key bindings
type keyMap struct {
	Quit   key.Binding
	Escape key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
	),
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
