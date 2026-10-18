package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/farseenmanekhan1232/tkeybr/internal/engine"
)

// UserData represents all stored user data
type UserData struct {
	Stats       *engine.Stats           `json:"stats"`
	Preferences *UserPreferences        `json:"preferences"`
	Progress    *engine.ProgressTracker `json:"progress"`
}

// UserPreferences represents user settings
type UserPreferences struct {
	KeyboardLayout    string `json:"keyboard_layout"`     // "qwerty", "dvorak", "colemak"
	ShowKeyboard      bool   `json:"show_keyboard"`
	WordMode          string `json:"word_mode"`           // "phonetic", "real", "mixed"
	Difficulty        int    `json:"difficulty"`
	WordCount         int    `json:"word_count"`
	Theme             string `json:"theme"`               // "minimal", "dracula", "nord", etc.
	ShowOnLeaderboard bool   `json:"show_on_leaderboard"` // Opt-in to public leaderboard
	DisplayName       string `json:"display_name"`        // Public display name for leaderboard
}

// LocalStorage handles local file storage
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage handler
func NewLocalStorage() *LocalStorage {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".tkeybr")
	return &LocalStorage{basePath: basePath}
}

// ensureDir ensures the storage directory exists
func (ls *LocalStorage) ensureDir() error {
	return os.MkdirAll(ls.basePath, 0700)
}

// SaveStats saves typing statistics
func (ls *LocalStorage) SaveStats(stats *engine.Stats) error {
	if err := ls.ensureDir(); err != nil {
		return err
	}

	path := filepath.Join(ls.basePath, "stats.json")
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// LoadStats loads typing statistics
func (ls *LocalStorage) LoadStats() (*engine.Stats, error) {
	path := filepath.Join(ls.basePath, "stats.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return engine.NewStats(), nil
		}
		return nil, fmt.Errorf("failed to read stats: %w", err)
	}

	var stats engine.Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	// Ensure KeyStats map is initialized
	if stats.KeyStats == nil {
		stats.KeyStats = make(map[rune]*engine.KeyStat)
	}

	return &stats, nil
}

// DefaultPreferences returns default user preferences
func DefaultPreferences() *UserPreferences {
	return &UserPreferences{
		KeyboardLayout: "qwerty",
		ShowKeyboard:   true,
		WordMode:       "real",
		Difficulty:     5,
		WordCount:      15,
		Theme:          "minimal",
	}
}

// SavePreferences saves user preferences
func (ls *LocalStorage) SavePreferences(prefs *UserPreferences) error {
	if err := ls.ensureDir(); err != nil {
		return err
	}

	path := filepath.Join(ls.basePath, "preferences.json")
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// LoadPreferences loads user preferences
func (ls *LocalStorage) LoadPreferences() (*UserPreferences, error) {
	path := filepath.Join(ls.basePath, "preferences.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPreferences(), nil
		}
		return nil, fmt.Errorf("failed to read preferences: %w", err)
	}

	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("failed to parse preferences: %w", err)
	}

	return &prefs, nil
}

// SaveProgress saves progress tracker data
func (ls *LocalStorage) SaveProgress(progress *engine.ProgressTracker) error {
	if err := ls.ensureDir(); err != nil {
		return err
	}

	path := filepath.Join(ls.basePath, "progress.json")
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// LoadProgress loads progress tracker data
func (ls *LocalStorage) LoadProgress(layoutName string) (*engine.ProgressTracker, error) {
	path := filepath.Join(ls.basePath, "progress.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return engine.NewProgressTracker(layoutName), nil
		}
		return nil, fmt.Errorf("failed to read progress: %w", err)
	}

	var progress engine.ProgressTracker
	if err := json.Unmarshal(data, &progress); err != nil {
		return nil, fmt.Errorf("failed to parse progress: %w", err)
	}

	// Restore the layout reference (not serialized)
	progress.SetLayout(progress.LayoutName)

	// Ensure Letters map is initialized
	if progress.Letters == nil {
		progress.Letters = make(map[rune]*engine.LetterProgress)
	}

	return &progress, nil
}

// SaveAll saves all user data
func (ls *LocalStorage) SaveAll(data *UserData) error {
	if data.Stats != nil {
		if err := ls.SaveStats(data.Stats); err != nil {
			return err
		}
	}
	if data.Preferences != nil {
		if err := ls.SavePreferences(data.Preferences); err != nil {
			return err
		}
	}
	if data.Progress != nil {
		if err := ls.SaveProgress(data.Progress); err != nil {
			return err
		}
	}
	return nil
}

// LoadAll loads all user data
func (ls *LocalStorage) LoadAll() (*UserData, error) {
	stats, err := ls.LoadStats()
	if err != nil {
		return nil, err
	}

	prefs, err := ls.LoadPreferences()
	if err != nil {
		return nil, err
	}

	progress, err := ls.LoadProgress(prefs.KeyboardLayout)
	if err != nil {
		return nil, err
	}

	return &UserData{
		Stats:       stats,
		Preferences: prefs,
		Progress:    progress,
	}, nil
}
