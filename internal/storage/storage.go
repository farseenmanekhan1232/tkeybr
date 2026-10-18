package storage

import (
	"github.com/farseenmanekhan1232/tkeybr/internal/engine"
)

// Storage defines the interface for data persistence
type Storage interface {
	// Stats operations
	SaveStats(stats *engine.Stats) error
	LoadStats() (*engine.Stats, error)

	// Preferences operations
	SavePreferences(prefs *UserPreferences) error
	LoadPreferences() (*UserPreferences, error)

	// Combined operations
	SaveAll(data *UserData) error
	LoadAll() (*UserData, error)
}

// SyncedStorage combines local and cloud storage with sync capability
type SyncedStorage struct {
	local  *LocalStorage
	cloud  *FirestoreStorage
	userID string
}

// NewSyncedStorage creates a new synced storage instance
func NewSyncedStorage(local *LocalStorage, cloud *FirestoreStorage, userID string) *SyncedStorage {
	return &SyncedStorage{
		local:  local,
		cloud:  cloud,
		userID: userID,
	}
}

// SaveStats saves stats to both local and cloud storage
func (s *SyncedStorage) SaveStats(stats *engine.Stats) error {
	// Always save locally first (fast, works offline)
	if err := s.local.SaveStats(stats); err != nil {
		return err
	}

	// Then sync to cloud (best effort, don't block on errors)
	if s.cloud != nil {
		go s.cloud.SaveStats(stats)
	}

	return nil
}

// LoadStats loads stats, preferring cloud if available
func (s *SyncedStorage) LoadStats() (*engine.Stats, error) {
	// Try loading from cloud first
	if s.cloud != nil {
		cloudStats, err := s.cloud.LoadStats()
		if err == nil && cloudStats != nil {
			// Also update local cache
			s.local.SaveStats(cloudStats)
			return cloudStats, nil
		}
	}

	// Fallback to local
	return s.local.LoadStats()
}

// SavePreferences saves preferences to both local and cloud storage
func (s *SyncedStorage) SavePreferences(prefs *UserPreferences) error {
	if err := s.local.SavePreferences(prefs); err != nil {
		return err
	}

	if s.cloud != nil {
		go s.cloud.SavePreferences(prefs)
	}

	return nil
}

// LoadPreferences loads preferences, preferring cloud if available
func (s *SyncedStorage) LoadPreferences() (*UserPreferences, error) {
	if s.cloud != nil {
		cloudPrefs, err := s.cloud.LoadPreferences()
		if err == nil && cloudPrefs != nil {
			s.local.SavePreferences(cloudPrefs)
			return cloudPrefs, nil
		}
	}

	return s.local.LoadPreferences()
}

// SaveProgress saves progress to both local and cloud storage
func (s *SyncedStorage) SaveProgress(progress *engine.ProgressTracker) error {
	if err := s.local.SaveProgress(progress); err != nil {
		return err
	}

	if s.cloud != nil {
		go s.cloud.SaveProgress(progress)
	}

	return nil
}

// LoadProgress loads progress, preferring cloud if available
func (s *SyncedStorage) LoadProgress(layoutName string) (*engine.ProgressTracker, error) {
	if s.cloud != nil {
		cloudProgress, err := s.cloud.LoadProgress(layoutName)
		if err == nil && cloudProgress != nil {
			s.local.SaveProgress(cloudProgress)
			return cloudProgress, nil
		}
	}

	return s.local.LoadProgress(layoutName)
}

// SaveAll saves all user data
func (s *SyncedStorage) SaveAll(data *UserData) error {
	if err := s.local.SaveAll(data); err != nil {
		return err
	}

	if s.cloud != nil {
		go s.cloud.SaveAll(data)
	}

	return nil
}

// LoadAll loads all user data
func (s *SyncedStorage) LoadAll() (*UserData, error) {
	if s.cloud != nil {
		cloudData, err := s.cloud.LoadAll()
		if err == nil && cloudData != nil {
			s.local.SaveAll(cloudData)
			return cloudData, nil
		}
	}

	return s.local.LoadAll()
}

// Sync synchronizes local and cloud data (merges stats)
func (s *SyncedStorage) Sync() error {
	if s.cloud == nil {
		return nil
	}

	// Load both
	localData, err := s.local.LoadAll()
	if err != nil {
		return err
	}

	cloudData, err := s.cloud.LoadAll()
	if err != nil {
		// Cloud not available, just push local
		return s.cloud.SaveAll(localData)
	}

	// Merge stats (take the maximum values)
	merged := mergeStats(localData.Stats, cloudData.Stats)

	// For preferences, prefer local (user's current session)
	merged.TotalCharacters = max(localData.Stats.TotalCharacters, cloudData.Stats.TotalCharacters)

	// Save merged data to both
	localData.Stats = merged
	if err := s.local.SaveStats(merged); err != nil {
		return err
	}
	if err := s.cloud.SaveStats(merged); err != nil {
		return err
	}

	return nil
}

// mergeStats merges two stats objects, taking maximum values
func mergeStats(a, b *engine.Stats) *engine.Stats {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	merged := engine.NewStats()

	// Take maximums for cumulative values
	merged.TotalCharacters = max(a.TotalCharacters, b.TotalCharacters)
	merged.TotalErrors = max(a.TotalErrors, b.TotalErrors)
	merged.TotalSessions = max(a.TotalSessions, b.TotalSessions)
	merged.TotalTimeSeconds = max(a.TotalTimeSeconds, b.TotalTimeSeconds)
	merged.BestWPM = max(a.BestWPM, b.BestWPM)
	merged.AverageWPM = max(a.AverageWPM, b.AverageWPM)

	// Merge key stats
	merged.KeyStats = make(map[rune]*engine.KeyStat)
	for k, v := range a.KeyStats {
		merged.KeyStats[k] = v
	}
	for k, v := range b.KeyStats {
		if existing, ok := merged.KeyStats[k]; ok {
			// Merge: take max counts
			existing.TotalTyped = max(existing.TotalTyped, v.TotalTyped)
			existing.ErrorCount = max(existing.ErrorCount, v.ErrorCount)
		} else {
			merged.KeyStats[k] = v
		}
	}

	return merged
}

func max[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Local returns the underlying local storage
func (s *SyncedStorage) Local() *LocalStorage {
	return s.local
}

// Cloud returns the underlying cloud storage (may be nil)
func (s *SyncedStorage) Cloud() *FirestoreStorage {
	return s.cloud
}
