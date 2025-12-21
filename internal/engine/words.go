package engine

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

// LetterProgress tracks progress for a single letter
type LetterProgress struct {
	Letter       rune      `json:"letter"`
	Unlocked     bool      `json:"unlocked"`
	TotalTyped   int       `json:"total_typed"`
	ErrorCount   int       `json:"error_count"`
	TotalTimeMs  int64     `json:"total_time_ms"`
	LastTyped    time.Time `json:"last_typed"`
	Confidence   float64   `json:"confidence"`
}

// ProgressTracker tracks overall typing progress and letter unlocking
type ProgressTracker struct {
	Layout           *KeyboardLayout           `json:"-"`
	LayoutName       string                    `json:"layout_name"`
	Letters          map[rune]*LetterProgress  `json:"letters"`
	UnlockedCount    int                       `json:"unlocked_count"`
	FocusLetter      rune                      `json:"focus_letter"`
	ConfidenceThresh float64                   `json:"confidence_threshold"`
	MinSamples       int                       `json:"min_samples"`
}

// Phonetic patterns using c=consonant, v=vowel
var phoneticPatterns = []string{
	"cvc",      // tin, set
	"cvcc",     // tent, mist  
	"ccvc",     // stop, trip
	"cvcv",     // data, silo
	"cvccv",    // pasta, haste
	"cvcvc",    // limit, river
	"ccvcc",    // stand, crisp
}

// NewProgressTracker creates a new progress tracker for a layout
func NewProgressTracker(layoutName string) *ProgressTracker {
	layout := GetLayout(layoutName)
	
	pt := &ProgressTracker{
		Layout:           layout,
		LayoutName:       layoutName,
		Letters:          make(map[rune]*LetterProgress),
		UnlockedCount:    4, // Start with first 4 letters (2 per hand)
		ConfidenceThresh: 0.7,
		MinSamples:       10,
	}

	// Initialize all letters based on layout's unlock order
	for i, letter := range layout.UnlockOrder {
		pt.Letters[letter] = &LetterProgress{
			Letter:   letter,
			Unlocked: i < pt.UnlockedCount,
		}
	}

	// Set initial focus letter
	if pt.UnlockedCount > 0 && pt.UnlockedCount <= len(layout.UnlockOrder) {
		pt.FocusLetter = layout.UnlockOrder[pt.UnlockedCount-1]
	}

	return pt
}

// SetLayout updates the layout (useful after loading from storage)
func (pt *ProgressTracker) SetLayout(layoutName string) {
	pt.Layout = GetLayout(layoutName)
	pt.LayoutName = layoutName
}

// GetUnlockedLetters returns currently unlocked letters
func (pt *ProgressTracker) GetUnlockedLetters() []rune {
	if pt.Layout == nil {
		pt.Layout = GetLayout(pt.LayoutName)
	}
	count := pt.UnlockedCount
	if count > len(pt.Layout.UnlockOrder) {
		count = len(pt.Layout.UnlockOrder)
	}
	return pt.Layout.UnlockOrder[:count]
}

// isVowel checks if a rune is a vowel
func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// GetUnlockedConsonants returns unlocked consonants
func (pt *ProgressTracker) GetUnlockedConsonants() []rune {
	var consonants []rune
	for _, l := range pt.GetUnlockedLetters() {
		if !isVowel(l) {
			consonants = append(consonants, l)
		}
	}
	return consonants
}

// GetUnlockedVowels returns unlocked vowels
func (pt *ProgressTracker) GetUnlockedVowels() []rune {
	var result []rune
	for _, l := range pt.GetUnlockedLetters() {
		if isVowel(l) {
			result = append(result, l)
		}
	}
	// Always need at least one vowel for word generation
	if len(result) == 0 {
		// Check if any vowels are in the unlock order and force unlock the first one
		for _, l := range pt.Layout.UnlockOrder {
			if isVowel(l) {
				result = append(result, l)
				break
			}
		}
	}
	if len(result) == 0 {
		return []rune{'e'} // Absolute fallback
	}
	return result
}

// RecordKeypress records a keypress for progress tracking
func (pt *ProgressTracker) RecordKeypress(letter rune, correct bool, timeMs int64) {
	letter = rune(strings.ToLower(string(letter))[0])
	
	progress, ok := pt.Letters[letter]
	if !ok {
		return
	}

	progress.TotalTyped++
	progress.TotalTimeMs += timeMs
	progress.LastTyped = time.Now()
	
	if !correct {
		progress.ErrorCount++
	}

	// Recalculate confidence
	pt.updateConfidence(letter)

	// Check if we should unlock next letter
	pt.checkUnlock()
}

// updateConfidence recalculates confidence for a letter
func (pt *ProgressTracker) updateConfidence(letter rune) {
	progress := pt.Letters[letter]
	if progress.TotalTyped < pt.MinSamples {
		progress.Confidence = 0
		return
	}

	// Calculate accuracy (0-1)
	accuracy := 1.0 - float64(progress.ErrorCount)/float64(progress.TotalTyped)
	
	// Calculate speed score (0-1)
	// Average time per keypress, target is ~200ms for good typing
	avgTimeMs := float64(progress.TotalTimeMs) / float64(progress.TotalTyped)
	speedScore := math.Min(1.0, 200.0/avgTimeMs)
	
	// Factor in key difficulty from layout
	keyDifficulty := 1.0
	if pt.Layout != nil {
		keyDifficulty = pt.Layout.GetKeyDifficulty(letter)
	}
	
	// Adjust speed expectation based on key difficulty
	// Harder keys are allowed to be slower
	adjustedSpeedScore := speedScore * keyDifficulty
	
	// Confidence is weighted: 60% accuracy, 40% speed
	progress.Confidence = math.Min(1.0, (accuracy * 0.6) + (adjustedSpeedScore * 0.4))
}

// checkUnlock checks if we should unlock the next letter
func (pt *ProgressTracker) checkUnlock() {
	if pt.UnlockedCount >= len(pt.Layout.UnlockOrder) {
		return // All letters unlocked
	}

	// Check if all current unlocked letters have sufficient confidence
	allConfident := true
	for _, letter := range pt.GetUnlockedLetters() {
		progress := pt.Letters[letter]
		if progress.TotalTyped < pt.MinSamples || progress.Confidence < pt.ConfidenceThresh {
			allConfident = false
			break
		}
	}

	if allConfident {
		// Unlock next letter
		pt.UnlockedCount++
		if pt.UnlockedCount <= len(pt.Layout.UnlockOrder) {
			newLetter := pt.Layout.UnlockOrder[pt.UnlockedCount-1]
			pt.Letters[newLetter].Unlocked = true
			pt.FocusLetter = newLetter
		}
	}
}

// GetWeakestLetter returns the letter with lowest confidence among unlocked
func (pt *ProgressTracker) GetWeakestLetter() rune {
	var weakest rune
	lowestConf := 2.0

	for _, letter := range pt.GetUnlockedLetters() {
		progress := pt.Letters[letter]
		if progress.TotalTyped >= pt.MinSamples && progress.Confidence < lowestConf {
			lowestConf = progress.Confidence
			weakest = letter
		}
	}

	if weakest == 0 {
		return pt.FocusLetter
	}
	return weakest
}

// GetLetterWeights returns weights for word generation (higher = more frequent in practice)
func (pt *ProgressTracker) GetLetterWeights() map[rune]float64 {
	weights := make(map[rune]float64)
	
	for _, letter := range pt.GetUnlockedLetters() {
		progress := pt.Letters[letter]
		
		// Base weight
		weight := 1.0
		
		// Lower confidence = higher weight (practice more)
		if progress.TotalTyped >= pt.MinSamples {
			// Invert confidence: low confidence -> high weight
			weight = 1.0 + (1.0 - progress.Confidence)
		} else {
			// New letters get medium-high weight to gather samples
			weight = 1.5
		}
		
		// Focus letter gets extra weight
		if letter == pt.FocusLetter {
			weight *= 1.5
		}
		
		// Weakest letter gets even more
		if letter == pt.GetWeakestLetter() {
			weight *= 1.3
		}
		
		weights[letter] = weight
	}
	
	return weights
}

// AdaptiveWordGenerator generates words using only unlocked letters
type AdaptiveWordGenerator struct {
	tracker *ProgressTracker
}

// NewAdaptiveWordGenerator creates a new adaptive word generator
func NewAdaptiveWordGenerator(tracker *ProgressTracker) *AdaptiveWordGenerator {
	return &AdaptiveWordGenerator{tracker: tracker}
}

// selectWeightedLetter selects a letter based on weights
func (g *AdaptiveWordGenerator) selectWeightedLetter(letters []rune, weights map[rune]float64) rune {
	if len(letters) == 0 {
		return 'e'
	}

	// Calculate total weight
	totalWeight := 0.0
	for _, l := range letters {
		w := weights[l]
		if w == 0 {
			w = 1.0
		}
		totalWeight += w
	}

	// Select based on weight
	r := rand.Float64() * totalWeight
	for _, l := range letters {
		w := weights[l]
		if w == 0 {
			w = 1.0
		}
		r -= w
		if r <= 0 {
			return l
		}
	}

	return letters[0]
}

// GenerateWord generates a single phonetic pseudo-word using unlocked letters
func (g *AdaptiveWordGenerator) GenerateWord() string {
	consonants := g.tracker.GetUnlockedConsonants()
	vowels := g.tracker.GetUnlockedVowels()
	weights := g.tracker.GetLetterWeights()

	// Need minimum letters
	if len(consonants) < 2 {
		// Use first few consonants from unlock order as fallback
		for _, l := range g.tracker.Layout.UnlockOrder {
			if !isVowel(l) && len(consonants) < 3 {
				consonants = append(consonants, l)
			}
		}
	}

	// Select pattern - simpler patterns when fewer letters unlocked
	maxPatternIdx := len(phoneticPatterns) - 1
	progressRatio := float64(g.tracker.UnlockedCount) / float64(len(g.tracker.Layout.UnlockOrder))
	maxPatternIdx = int(float64(len(phoneticPatterns)) * progressRatio)
	if maxPatternIdx < 1 {
		maxPatternIdx = 1
	}
	if maxPatternIdx >= len(phoneticPatterns) {
		maxPatternIdx = len(phoneticPatterns) - 1
	}
	
	pattern := phoneticPatterns[rand.Intn(maxPatternIdx+1)]

	var word strings.Builder
	var lastChar rune
	
	for _, p := range pattern {
		var char rune
		if p == 'c' {
			char = g.selectWeightedLetter(consonants, weights)
		} else {
			char = g.selectWeightedLetter(vowels, weights)
		}
		
		// Avoid same-finger bigrams when possible (reroll once)
		if lastChar != 0 && g.tracker.Layout.SameFinger(lastChar, char) {
			if p == 'c' {
				char = g.selectWeightedLetter(consonants, weights)
			} else {
				char = g.selectWeightedLetter(vowels, weights)
			}
		}
		
		word.WriteRune(char)
		lastChar = char
	}

	return word.String()
}

// GenerateText generates practice text with specified word count
func (g *AdaptiveWordGenerator) GenerateText(wordCount int) string {
	words := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		words[i] = g.GenerateWord()
	}
	return strings.Join(words, " ")
}

// GetProgressSummary returns a summary of current progress
func (pt *ProgressTracker) GetProgressSummary() ProgressSummary {
	unlocked := pt.GetUnlockedLetters()
	return ProgressSummary{
		UnlockedCount: pt.UnlockedCount,
		TotalLetters:  len(pt.Layout.UnlockOrder),
		FocusLetter:   pt.FocusLetter,
		WeakestLetter: pt.GetWeakestLetter(),
		Unlocked:      unlocked,
		LayoutName:    pt.LayoutName,
	}
}

// ProgressSummary is a JSON-friendly progress summary
type ProgressSummary struct {
	UnlockedCount int    `json:"unlocked_count"`
	TotalLetters  int    `json:"total_letters"`
	FocusLetter   rune   `json:"focus_letter"`
	WeakestLetter rune   `json:"weakest_letter"`
	Unlocked      []rune `json:"unlocked"`
	LayoutName    string `json:"layout_name"`
}

// ResetProgress resets all progress (keeps layout)
func (pt *ProgressTracker) ResetProgress() {
	pt.UnlockedCount = 4
	for _, letter := range pt.Layout.UnlockOrder {
		pt.Letters[letter] = &LetterProgress{
			Letter:   letter,
			Unlocked: false,
		}
	}
	for i := 0; i < pt.UnlockedCount; i++ {
		pt.Letters[pt.Layout.UnlockOrder[i]].Unlocked = true
	}
	if pt.UnlockedCount > 0 {
		pt.FocusLetter = pt.Layout.UnlockOrder[pt.UnlockedCount-1]
	}
}
