package engine

// KeyboardLayout defines a keyboard layout with finger mappings and difficulty
type KeyboardLayout struct {
	Name        string
	Rows        [][]rune          // Physical key rows
	FingerMap   map[rune]Finger   // Which finger types each key
	RowMap      map[rune]int      // Which row each key is on (0=home, -1=top, 1=bottom)
	UnlockOrder []rune            // Recommended order to learn letters
}

// Finger represents which finger types a key
type Finger int

const (
	LeftPinky Finger = iota
	LeftRing
	LeftMiddle
	LeftIndex
	RightIndex
	RightMiddle
	RightRing
	RightPinky
	Thumb
)

// FingerDifficulty - lower is easier (index/middle are strongest)
var FingerDifficulty = map[Finger]float64{
	LeftPinky:   1.4,  // Weakest
	LeftRing:    1.2,
	LeftMiddle:  1.0,  // Strong
	LeftIndex:   0.9,  // Strongest
	RightIndex:  0.9,
	RightMiddle: 1.0,
	RightRing:   1.2,
	RightPinky:  1.4,
	Thumb:       0.8,  // Spacebar
}

// RowDifficulty - home row is easiest
var RowDifficulty = map[int]float64{
	0:  1.0, // Home row - easiest
	-1: 1.2, // Top row
	1:  1.3, // Bottom row - hardest (reach down)
}

// GetKeyDifficulty returns the difficulty score for a key (lower = easier)
func (kl *KeyboardLayout) GetKeyDifficulty(key rune) float64 {
	finger, ok := kl.FingerMap[key]
	if !ok {
		return 1.5 // Unknown key
	}
	row, ok := kl.RowMap[key]
	if !ok {
		return 1.5
	}
	
	return FingerDifficulty[finger] * RowDifficulty[row]
}

// QWERTY Layout
var QWERTYLayout = &KeyboardLayout{
	Name: "QWERTY",
	Rows: [][]rune{
		{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},     // Top row
		{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'},          // Home row
		{'z', 'x', 'c', 'v', 'b', 'n', 'm'},                     // Bottom row
	},
	FingerMap: map[rune]Finger{
		// Top row
		'q': LeftPinky, 'w': LeftRing, 'e': LeftMiddle, 'r': LeftIndex, 't': LeftIndex,
		'y': RightIndex, 'u': RightIndex, 'i': RightMiddle, 'o': RightRing, 'p': RightPinky,
		// Home row
		'a': LeftPinky, 's': LeftRing, 'd': LeftMiddle, 'f': LeftIndex, 'g': LeftIndex,
		'h': RightIndex, 'j': RightIndex, 'k': RightMiddle, 'l': RightRing,
		// Bottom row
		'z': LeftPinky, 'x': LeftRing, 'c': LeftMiddle, 'v': LeftIndex, 'b': LeftIndex,
		'n': RightIndex, 'm': RightIndex,
		// Space
		' ': Thumb,
	},
	RowMap: map[rune]int{
		// Top row (-1)
		'q': -1, 'w': -1, 'e': -1, 'r': -1, 't': -1,
		'y': -1, 'u': -1, 'i': -1, 'o': -1, 'p': -1,
		// Home row (0)
		'a': 0, 's': 0, 'd': 0, 'f': 0, 'g': 0,
		'h': 0, 'j': 0, 'k': 0, 'l': 0,
		// Bottom row (1)
		'z': 1, 'x': 1, 'c': 1, 'v': 1, 'b': 1, 'n': 1, 'm': 1,
		// Space
		' ': 0,
	},
	// Unlock order: Home row first (strong fingers), then expand outward
	// Prioritizes: 1) Home row, 2) Strong fingers (index/middle), 3) Frequency
	UnlockOrder: []rune{
		// Phase 1: Home row, index fingers (easiest)
		'f', 'j',
		// Phase 2: Home row, middle fingers
		'd', 'k',
		// Phase 3: Home row, ring fingers  
		's', 'l',
		// Phase 4: Home row extend + top row index
		'g', 'h', 'r', 'u',
		// Phase 5: Top row middle
		'e', 'i',
		// Phase 6: Top row ring + common
		'w', 'o', 't', 'y',
		// Phase 7: Bottom row index/middle (common letters)
		'n', 'm', 'v', 'c',
		// Phase 8: Home row pinky
		'a',
		// Phase 9: Top row pinky + bottom
		'q', 'p', 'b',
		// Phase 10: Difficult letters
		'x', 'z',
	},
}

// Dvorak Layout  
var DvorakLayout = &KeyboardLayout{
	Name: "Dvorak",
	Rows: [][]rune{
		{'\'', ',', '.', 'p', 'y', 'f', 'g', 'c', 'r', 'l'},
		{'a', 'o', 'e', 'u', 'i', 'd', 'h', 't', 'n', 's'},
		{';', 'q', 'j', 'k', 'x', 'b', 'm', 'w', 'v', 'z'},
	},
	FingerMap: map[rune]Finger{
		// Top row
		'\'': LeftPinky, ',': LeftRing, '.': LeftMiddle, 'p': LeftIndex, 'y': LeftIndex,
		'f': RightIndex, 'g': RightIndex, 'c': RightMiddle, 'r': RightRing, 'l': RightPinky,
		// Home row - Most common letters are here in Dvorak!
		'a': LeftPinky, 'o': LeftRing, 'e': LeftMiddle, 'u': LeftIndex, 'i': LeftIndex,
		'd': RightIndex, 'h': RightIndex, 't': RightMiddle, 'n': RightRing, 's': RightPinky,
		// Bottom row
		';': LeftPinky, 'q': LeftRing, 'j': LeftMiddle, 'k': LeftIndex, 'x': LeftIndex,
		'b': RightIndex, 'm': RightIndex, 'w': RightMiddle, 'v': RightRing, 'z': RightPinky,
		' ': Thumb,
	},
	RowMap: map[rune]int{
		'\'': -1, ',': -1, '.': -1, 'p': -1, 'y': -1,
		'f': -1, 'g': -1, 'c': -1, 'r': -1, 'l': -1,
		'a': 0, 'o': 0, 'e': 0, 'u': 0, 'i': 0,
		'd': 0, 'h': 0, 't': 0, 'n': 0, 's': 0,
		';': 1, 'q': 1, 'j': 1, 'k': 1, 'x': 1,
		'b': 1, 'm': 1, 'w': 1, 'v': 1, 'z': 1,
		' ': 0,
	},
	// Dvorak: Home row has vowels + common consonants
	UnlockOrder: []rune{
		// Home row - all on home in Dvorak!
		'u', 'i', 'd', 'h',     // Index fingers
		'e', 't',               // Middle fingers  
		'o', 'n',               // Ring fingers
		'a', 's',               // Pinky (but common letters)
		// Top/bottom common
		'r', 'l', 'c', 'g',
		'p', 'y', 'f',
		'm', 'w', 'b',
		'k', 'v', 'j', 'q', 'x', 'z',
	},
}

// Colemak Layout
var ColemakLayout = &KeyboardLayout{
	Name: "Colemak",
	Rows: [][]rune{
		{'q', 'w', 'f', 'p', 'g', 'j', 'l', 'u', 'y', ';'},
		{'a', 'r', 's', 't', 'd', 'h', 'n', 'e', 'i', 'o'},
		{'z', 'x', 'c', 'v', 'b', 'k', 'm'},
	},
	FingerMap: map[rune]Finger{
		// Top row
		'q': LeftPinky, 'w': LeftRing, 'f': LeftMiddle, 'p': LeftIndex, 'g': LeftIndex,
		'j': RightIndex, 'l': RightIndex, 'u': RightMiddle, 'y': RightRing, ';': RightPinky,
		// Home row
		'a': LeftPinky, 'r': LeftRing, 's': LeftMiddle, 't': LeftIndex, 'd': LeftIndex,
		'h': RightIndex, 'n': RightIndex, 'e': RightMiddle, 'i': RightRing, 'o': RightPinky,
		// Bottom row
		'z': LeftPinky, 'x': LeftRing, 'c': LeftMiddle, 'v': LeftIndex, 'b': LeftIndex,
		'k': RightIndex, 'm': RightIndex,
		' ': Thumb,
	},
	RowMap: map[rune]int{
		'q': -1, 'w': -1, 'f': -1, 'p': -1, 'g': -1,
		'j': -1, 'l': -1, 'u': -1, 'y': -1, ';': -1,
		'a': 0, 'r': 0, 's': 0, 't': 0, 'd': 0,
		'h': 0, 'n': 0, 'e': 0, 'i': 0, 'o': 0,
		'z': 1, 'x': 1, 'c': 1, 'v': 1, 'b': 1, 'k': 1, 'm': 1,
		' ': 0,
	},
	// Colemak: Optimized home row
	UnlockOrder: []rune{
		// Home row - common letters
		't', 'd', 'h', 'n',     // Index fingers
		's', 'e',               // Middle fingers
		'r', 'i',               // Ring fingers  
		'a', 'o',               // Pinky
		// Expand out
		'l', 'u', 'f', 'p',
		'g', 'j', 'w', 'y',
		'c', 'm', 'v', 'b', 'k',
		'q', 'x', 'z',
	},
}

// GetLayout returns a layout by name
func GetLayout(name string) *KeyboardLayout {
	switch name {
	case "dvorak":
		return DvorakLayout
	case "colemak":
		return ColemakLayout
	default:
		return QWERTYLayout
	}
}

// GetAllLayouts returns all available layouts
func GetAllLayouts() []*KeyboardLayout {
	return []*KeyboardLayout{QWERTYLayout, DvorakLayout, ColemakLayout}
}

// SameFinger checks if two keys are typed with the same finger
func (kl *KeyboardLayout) SameFinger(a, b rune) bool {
	fingerA, okA := kl.FingerMap[a]
	fingerB, okB := kl.FingerMap[b]
	if !okA || !okB {
		return false
	}
	return fingerA == fingerB
}

// GetBigramDifficulty returns difficulty of typing two keys in sequence
func (kl *KeyboardLayout) GetBigramDifficulty(a, b rune) float64 {
	diffA := kl.GetKeyDifficulty(a)
	diffB := kl.GetKeyDifficulty(b)
	base := (diffA + diffB) / 2

	// Same finger bigrams are harder
	if kl.SameFinger(a, b) {
		base *= 1.5
	}

	return base
}
