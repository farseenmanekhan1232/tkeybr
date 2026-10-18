package engine

import (
	"time"
)

// TypingState represents the current state of a typing session
type TypingState struct {
	TargetText    string        // The text to type
	TypedText     string        // What the user has typed
	Position      int           // Current cursor position
	Errors        []int         // Positions where errors occurred
	StartTime     time.Time     // When typing started
	EndTime       time.Time     // When typing ended (if complete)
	IsComplete    bool          // Whether the text is fully typed
	CharTimes     []time.Time   // Timestamp for each character typed
	CharCorrect   []bool        // Whether each character was correct
}

// NewTypingState creates a new typing state for the given target text
func NewTypingState(targetText string) *TypingState {
	return &TypingState{
		TargetText:  targetText,
		TypedText:   "",
		Position:    0,
		Errors:      make([]int, 0),
		CharTimes:   make([]time.Time, 0),
		CharCorrect: make([]bool, 0),
	}
}

// TypeChar processes a typed character
func (ts *TypingState) TypeChar(c rune) {
	if ts.IsComplete {
		return
	}

	// Start timer on first character
	if ts.Position == 0 && ts.StartTime.IsZero() {
		ts.StartTime = time.Now()
	}

	if ts.Position >= len(ts.TargetText) {
		return
	}

	// Record timing
	ts.CharTimes = append(ts.CharTimes, time.Now())

	// Check if correct
	targetChar := rune(ts.TargetText[ts.Position])
	isCorrect := c == targetChar
	ts.CharCorrect = append(ts.CharCorrect, isCorrect)

	if !isCorrect {
		ts.Errors = append(ts.Errors, ts.Position)
	}

	ts.TypedText += string(c)
	ts.Position++

	// Check if complete
	if ts.Position >= len(ts.TargetText) {
		ts.IsComplete = true
		ts.EndTime = time.Now()
	}
}

// Backspace removes the last typed character
func (ts *TypingState) Backspace() {
	if ts.Position > 0 {
		ts.Position--
		ts.TypedText = ts.TypedText[:len(ts.TypedText)-1]
		
		// Remove from tracking
		if len(ts.CharTimes) > 0 {
			ts.CharTimes = ts.CharTimes[:len(ts.CharTimes)-1]
		}
		if len(ts.CharCorrect) > 0 {
			ts.CharCorrect = ts.CharCorrect[:len(ts.CharCorrect)-1]
		}
		
		// Remove error if it was at this position
		for i := len(ts.Errors) - 1; i >= 0; i-- {
			if ts.Errors[i] == ts.Position {
				ts.Errors = append(ts.Errors[:i], ts.Errors[i+1:]...)
			}
		}
		
		ts.IsComplete = false
	}
}

// GetCurrentChar returns the current character to type
func (ts *TypingState) GetCurrentChar() rune {
	if ts.Position < len(ts.TargetText) {
		return rune(ts.TargetText[ts.Position])
	}
	return 0
}

// GetElapsedTime returns the elapsed time since typing started
func (ts *TypingState) GetElapsedTime() time.Duration {
	if ts.StartTime.IsZero() {
		return 0
	}
	if ts.IsComplete {
		return ts.EndTime.Sub(ts.StartTime)
	}
	return time.Since(ts.StartTime)
}

// GetCharStatus returns the status of a character at a position
// Returns: "correct", "incorrect", "current", "pending"
func (ts *TypingState) GetCharStatus(pos int) string {
	if pos < ts.Position {
		// Already typed
		for _, errPos := range ts.Errors {
			if errPos == pos {
				return "incorrect"
			}
		}
		return "correct"
	} else if pos == ts.Position {
		return "current"
	}
	return "pending"
}

// Reset resets the typing state for a new session
func (ts *TypingState) Reset() {
	ts.TypedText = ""
	ts.Position = 0
	ts.Errors = make([]int, 0)
	ts.StartTime = time.Time{}
	ts.EndTime = time.Time{}
	ts.IsComplete = false
	ts.CharTimes = make([]time.Time, 0)
	ts.CharCorrect = make([]bool, 0)
}

// SetNewText sets new target text and resets state
func (ts *TypingState) SetNewText(text string) {
	ts.TargetText = text
	ts.Reset()
}
