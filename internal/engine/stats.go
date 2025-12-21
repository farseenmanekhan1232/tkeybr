package engine

import (
	"time"
)

// Stats holds typing statistics
type Stats struct {
	// Current session
	CurrentWPM      float64
	CurrentCPM      float64
	CurrentAccuracy float64
	CharactersTyped int
	ErrorCount      int
	
	// Historical data
	TotalCharacters  int
	TotalErrors      int
	TotalSessions    int
	TotalTimeSeconds float64
	BestWPM          float64
	AverageWPM       float64
	
	// Per-key statistics
	KeyStats map[rune]*KeyStat
}

// KeyStat holds statistics for a single key
type KeyStat struct {
	TotalTyped    int
	ErrorCount    int
	AverageTimeMs float64 // Average time to type this key
	LastTyped     time.Time
}

// NewStats creates a new Stats instance
func NewStats() *Stats {
	return &Stats{
		KeyStats: make(map[rune]*KeyStat),
	}
}

// CalculateWPM calculates words per minute
// Standard: 1 word = 5 characters
func CalculateWPM(characters int, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	words := float64(characters) / 5.0
	minutes := duration.Minutes()
	if minutes == 0 {
		return 0
	}
	return words / minutes
}

// CalculateCPM calculates characters per minute
func CalculateCPM(characters int, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	minutes := duration.Minutes()
	if minutes == 0 {
		return 0
	}
	return float64(characters) / minutes
}

// CalculateAccuracy calculates accuracy percentage
func CalculateAccuracy(totalChars, errors int) float64 {
	if totalChars == 0 {
		return 100.0
	}
	return float64(totalChars-errors) / float64(totalChars) * 100
}

// UpdateFromTypingState updates stats from a typing state
func (s *Stats) UpdateFromTypingState(ts *TypingState) {
	duration := ts.GetElapsedTime()
	
	s.CharactersTyped = ts.Position
	s.ErrorCount = len(ts.Errors)
	
	// Calculate correct characters (total - errors)
	correctChars := ts.Position - len(ts.Errors)
	if correctChars < 0 {
		correctChars = 0
	}
	
	// WPM is based on CORRECT characters only (not raw typing)
	s.CurrentWPM = CalculateWPM(correctChars, duration)
	s.CurrentCPM = CalculateCPM(correctChars, duration)
	s.CurrentAccuracy = CalculateAccuracy(ts.Position, len(ts.Errors))

	// Update per-key stats
	for i, correct := range ts.CharCorrect {
		if i < len(ts.TargetText) {
			char := rune(ts.TargetText[i])
			s.updateKeyStats(char, correct)
		}
	}
}

// updateKeyStats updates statistics for a single key
func (s *Stats) updateKeyStats(key rune, correct bool) {
	if _, ok := s.KeyStats[key]; !ok {
		s.KeyStats[key] = &KeyStat{}
	}
	
	stat := s.KeyStats[key]
	stat.TotalTyped++
	if !correct {
		stat.ErrorCount++
	}
	stat.LastTyped = time.Now()
}

// CompleteSession marks a session as complete and updates historical stats
func (s *Stats) CompleteSession(ts *TypingState) {
	s.UpdateFromTypingState(ts)
	
	// Track correct characters, not total
	correctChars := ts.Position - len(ts.Errors)
	if correctChars < 0 {
		correctChars = 0
	}
	
	s.TotalCharacters += correctChars // Only count correct chars
	s.TotalErrors += len(ts.Errors)
	s.TotalSessions++
	s.TotalTimeSeconds += ts.GetElapsedTime().Seconds()
	
	if s.CurrentWPM > s.BestWPM {
		s.BestWPM = s.CurrentWPM
	}
	
	// Recalculate average WPM based on correct characters
	if s.TotalTimeSeconds > 0 {
		totalWords := float64(s.TotalCharacters) / 5.0
		totalMinutes := s.TotalTimeSeconds / 60.0
		if totalMinutes > 0 {
			s.AverageWPM = totalWords / totalMinutes
		}
	}
}

// GetWeakKeys returns keys with high error rates (for adaptive learning)
func (s *Stats) GetWeakKeys(threshold float64) map[rune]float64 {
	weakKeys := make(map[rune]float64)
	
	for key, stat := range s.KeyStats {
		if stat.TotalTyped >= 5 { // Need enough data
			errorRate := float64(stat.ErrorCount) / float64(stat.TotalTyped)
			if errorRate > threshold {
				weakKeys[key] = errorRate
			}
		}
	}
	
	return weakKeys
}

// GetKeyErrorRate returns the error rate for a specific key
func (s *Stats) GetKeyErrorRate(key rune) float64 {
	if stat, ok := s.KeyStats[key]; ok && stat.TotalTyped > 0 {
		return float64(stat.ErrorCount) / float64(stat.TotalTyped)
	}
	return 0
}
