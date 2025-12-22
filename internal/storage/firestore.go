package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/farseenmanekhan1232/tkeybr/internal/engine"
)

// FirestoreStorage handles cloud storage using Firestore
type FirestoreStorage struct {
	client *firestore.Client
	userID string
	ctx    context.Context
}

// NewFirestoreStorage creates a new Firestore storage handler
func NewFirestoreStorage(client *firestore.Client, userID string) *FirestoreStorage {
	return &FirestoreStorage{
		client: client,
		userID: userID,
		ctx:    context.Background(),
	}
}

// userDoc returns the user's document reference
func (fs *FirestoreStorage) userDoc() *firestore.DocumentRef {
	return fs.client.Collection("users").Doc(fs.userID)
}

// firestoreStats is the Firestore representation of Stats
type firestoreStats struct {
	TotalCharacters  int                        `firestore:"totalCharacters"`
	TotalErrors      int                        `firestore:"totalErrors"`
	TotalSessions    int                        `firestore:"totalSessions"`
	TotalTimeSeconds float64                    `firestore:"totalTimeSeconds"`
	BestWPM          float64                    `firestore:"bestWPM"`
	AverageWPM       float64                    `firestore:"averageWPM"`
	KeyStats         map[string]firestoreKeyStat `firestore:"keyStats"`
	UpdatedAt        time.Time                  `firestore:"updatedAt"`
}

// firestoreKeyStat is the Firestore representation of KeyStat
type firestoreKeyStat struct {
	TotalTyped    int       `firestore:"totalTyped"`
	ErrorCount    int       `firestore:"errorCount"`
	AverageTimeMs float64   `firestore:"averageTimeMs"`
	LastTyped     time.Time `firestore:"lastTyped"`
}

// firestorePreferences is the Firestore representation of UserPreferences
type firestorePreferences struct {
	KeyboardLayout string    `firestore:"keyboardLayout"`
	ShowKeyboard   bool      `firestore:"showKeyboard"`
	WordMode       string    `firestore:"wordMode"`
	Difficulty     int       `firestore:"difficulty"`
	WordCount      int       `firestore:"wordCount"`
	Theme          string    `firestore:"theme"`
	UpdatedAt      time.Time `firestore:"updatedAt"`
}

// SaveStats saves typing statistics to Firestore
func (fs *FirestoreStorage) SaveStats(stats *engine.Stats) error {
	if stats == nil {
		return nil
	}

	// Convert KeyStats map (rune keys to string keys for Firestore)
	keyStats := make(map[string]firestoreKeyStat)
	for k, v := range stats.KeyStats {
		keyStats[string(k)] = firestoreKeyStat{
			TotalTyped:    v.TotalTyped,
			ErrorCount:    v.ErrorCount,
			AverageTimeMs: v.AverageTimeMs,
			LastTyped:     v.LastTyped,
		}
	}

	fsStats := firestoreStats{
		TotalCharacters:  stats.TotalCharacters,
		TotalErrors:      stats.TotalErrors,
		TotalSessions:    stats.TotalSessions,
		TotalTimeSeconds: stats.TotalTimeSeconds,
		BestWPM:          stats.BestWPM,
		AverageWPM:       stats.AverageWPM,
		KeyStats:         keyStats,
		UpdatedAt:        time.Now(),
	}

	_, err := fs.userDoc().Collection("data").Doc("stats").Set(fs.ctx, fsStats)
	if err != nil {
		return fmt.Errorf("failed to save stats to Firestore: %w", err)
	}

	return nil
}

// LoadStats loads typing statistics from Firestore
func (fs *FirestoreStorage) LoadStats() (*engine.Stats, error) {
	doc, err := fs.userDoc().Collection("data").Doc("stats").Get(fs.ctx)
	if err != nil {
		// Document doesn't exist, return new stats
		return engine.NewStats(), nil
	}

	var fsStats firestoreStats
	if err := doc.DataTo(&fsStats); err != nil {
		return nil, fmt.Errorf("failed to parse stats from Firestore: %w", err)
	}

	// Convert back to engine.Stats
	stats := &engine.Stats{
		TotalCharacters:  fsStats.TotalCharacters,
		TotalErrors:      fsStats.TotalErrors,
		TotalSessions:    fsStats.TotalSessions,
		TotalTimeSeconds: fsStats.TotalTimeSeconds,
		BestWPM:          fsStats.BestWPM,
		AverageWPM:       fsStats.AverageWPM,
		KeyStats:         make(map[rune]*engine.KeyStat),
	}

	// Convert string keys back to rune keys
	for k, v := range fsStats.KeyStats {
		if len(k) > 0 {
			runes := []rune(k)
			stats.KeyStats[runes[0]] = &engine.KeyStat{
				TotalTyped:    v.TotalTyped,
				ErrorCount:    v.ErrorCount,
				AverageTimeMs: v.AverageTimeMs,
				LastTyped:     v.LastTyped,
			}
		}
	}

	return stats, nil
}

// SavePreferences saves user preferences to Firestore
func (fs *FirestoreStorage) SavePreferences(prefs *UserPreferences) error {
	if prefs == nil {
		return nil
	}

	fsPrefs := firestorePreferences{
		KeyboardLayout: prefs.KeyboardLayout,
		ShowKeyboard:   prefs.ShowKeyboard,
		WordMode:       prefs.WordMode,
		Difficulty:     prefs.Difficulty,
		WordCount:      prefs.WordCount,
		Theme:          prefs.Theme,
		UpdatedAt:      time.Now(),
	}

	_, err := fs.userDoc().Collection("data").Doc("preferences").Set(fs.ctx, fsPrefs)
	if err != nil {
		return fmt.Errorf("failed to save preferences to Firestore: %w", err)
	}

	return nil
}

// LoadPreferences loads user preferences from Firestore
func (fs *FirestoreStorage) LoadPreferences() (*UserPreferences, error) {
	doc, err := fs.userDoc().Collection("data").Doc("preferences").Get(fs.ctx)
	if err != nil {
		return DefaultPreferences(), nil
	}

	var fsPrefs firestorePreferences
	if err := doc.DataTo(&fsPrefs); err != nil {
		return nil, fmt.Errorf("failed to parse preferences from Firestore: %w", err)
	}

	return &UserPreferences{
		KeyboardLayout: fsPrefs.KeyboardLayout,
		ShowKeyboard:   fsPrefs.ShowKeyboard,
		WordMode:       fsPrefs.WordMode,
		Difficulty:     fsPrefs.Difficulty,
		WordCount:      fsPrefs.WordCount,
		Theme:          fsPrefs.Theme,
	}, nil
}

// SaveAll saves all user data to Firestore
func (fs *FirestoreStorage) SaveAll(data *UserData) error {
	if data.Stats != nil {
		if err := fs.SaveStats(data.Stats); err != nil {
			return err
		}
	}
	if data.Preferences != nil {
		if err := fs.SavePreferences(data.Preferences); err != nil {
			return err
		}
	}
	return nil
}

// LoadAll loads all user data from Firestore
func (fs *FirestoreStorage) LoadAll() (*UserData, error) {
	stats, err := fs.LoadStats()
	if err != nil {
		return nil, err
	}

	prefs, err := fs.LoadPreferences()
	if err != nil {
		return nil, err
	}

	return &UserData{
		Stats:       stats,
		Preferences: prefs,
	}, nil
}

// CreateUserProfile creates or updates the user profile in Firestore
func (fs *FirestoreStorage) CreateUserProfile(email, name string) error {
	profile := map[string]interface{}{
		"email":     email,
		"name":      name,
		"createdAt": firestore.ServerTimestamp,
		"updatedAt": firestore.ServerTimestamp,
	}

	_, err := fs.userDoc().Set(fs.ctx, profile, firestore.MergeAll)
	return err
}
