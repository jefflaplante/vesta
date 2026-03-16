package api

import (
	"testing"
)

func TestCharToCode_Letters(t *testing.T) {
	tests := []struct {
		char rune
		want int
	}{
		{'A', 1}, {'B', 2}, {'C', 3}, {'D', 4}, {'E', 5},
		{'F', 6}, {'G', 7}, {'H', 8}, {'I', 9}, {'J', 10},
		{'K', 11}, {'L', 12}, {'M', 13}, {'N', 14}, {'O', 15},
		{'P', 16}, {'Q', 17}, {'R', 18}, {'S', 19}, {'T', 20},
		{'U', 21}, {'V', 22}, {'W', 23}, {'X', 24}, {'Y', 25}, {'Z', 26},
	}

	for _, tt := range tests {
		got := charToCode(tt.char)
		if got != tt.want {
			t.Errorf("charToCode(%c) = %d, want %d", tt.char, got, tt.want)
		}
	}
}

func TestCharToCode_LowercaseConvertsToUpper(t *testing.T) {
	tests := []struct {
		char rune
		want int
	}{
		{'a', 1}, {'b', 2}, {'c', 3}, {'z', 26},
	}

	for _, tt := range tests {
		got := charToCode(tt.char)
		if got != tt.want {
			t.Errorf("charToCode(%c) = %d, want %d", tt.char, got, tt.want)
		}
	}
}

func TestCharToCode_Digits(t *testing.T) {
	tests := []struct {
		char rune
		want int
	}{
		{'1', 27}, {'2', 28}, {'3', 29}, {'4', 30}, {'5', 31},
		{'6', 32}, {'7', 33}, {'8', 34}, {'9', 35}, {'0', 36},
	}

	for _, tt := range tests {
		got := charToCode(tt.char)
		if got != tt.want {
			t.Errorf("charToCode(%c) = %d, want %d", tt.char, got, tt.want)
		}
	}
}

func TestCharToCode_SpecialChars(t *testing.T) {
	tests := []struct {
		char rune
		want int
	}{
		{'!', 37}, {'@', 38}, {'#', 39}, {'$', 40},
		{'(', 41}, {')', 42}, {'*', 43}, {'-', 44},
		{'+', 46}, {'&', 47}, {'=', 48}, {';', 49}, {':', 50},
		{'\'', 52}, {'"', 53}, {'%', 54}, {',', 55},
		{'.', 56}, {'/', 59}, {'?', 60},
	}

	for _, tt := range tests {
		got := charToCode(tt.char)
		if got != tt.want {
			t.Errorf("charToCode(%c) = %d, want %d", tt.char, got, tt.want)
		}
	}
}

func TestCharToCode_Space(t *testing.T) {
	got := charToCode(' ')
	if got != 0 {
		t.Errorf("charToCode(' ') = %d, want 0", got)
	}
}

func TestCharToCode_UnknownCharsReturnZero(t *testing.T) {
	unknowns := []rune{'~', '`', '[', ']', '{', '}', '\\', '|', '^'}
	for _, r := range unknowns {
		got := charToCode(r)
		if got != 0 {
			t.Errorf("charToCode(%c) = %d, want 0 for unknown char", r, got)
		}
	}
}

func TestParseEscapeCode_Colors(t *testing.T) {
	tests := []struct {
		input string
		want  int
		ok    bool
	}{
		{"red", 63, true},
		{"orange", 64, true},
		{"yellow", 65, true},
		{"green", 66, true},
		{"blue", 67, true},
		{"violet", 68, true},
		{"white", 69, true},
		{"black", 70, true},
		{"filled", 71, true},
		{"RED", 63, true}, // case insensitive
		{"Red", 63, true}, // case insensitive
		{"deg", 62, true}, // degree
		{"<3", 62, true},  // heart
	}

	for _, tt := range tests {
		got, ok := parseEscapeCode(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("parseEscapeCode(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseEscapeCode_Numeric(t *testing.T) {
	tests := []struct {
		input string
		want  int
		ok    bool
	}{
		{"0", 0, true},
		{"1", 1, true},
		{"63", 63, true},
		{"71", 71, true},
	}

	for _, tt := range tests {
		got, ok := parseEscapeCode(tt.input)
		if ok != tt.ok || (ok && got != tt.want) {
			t.Errorf("parseEscapeCode(%q) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestFormat_BasicMessage(t *testing.T) {
	result := Format("HELLO", "note", false)

	if len(result.Characters) != 3 {
		t.Errorf("Format returned %d rows, want 3 for note", len(result.Characters))
	}
	if len(result.Characters[0]) != 15 {
		t.Errorf("Format returned %d cols, want 15 for note", len(result.Characters[0]))
	}

	// H=8, E=5, L=12, L=12, O=15
	expected := []int{8, 5, 12, 12, 15}
	for i, want := range expected {
		if result.Characters[0][i] != want {
			t.Errorf("result.Characters[0][%d] = %d, want %d", i, result.Characters[0][i], want)
		}
	}
}

func TestFormat_Multiline(t *testing.T) {
	result := Format("A\\nB\\nC", "note", false)

	if result.Characters[0][0] != 1 { // A
		t.Errorf("Row 0 first char = %d, want 1 (A)", result.Characters[0][0])
	}
	if result.Characters[1][0] != 2 { // B
		t.Errorf("Row 1 first char = %d, want 2 (B)", result.Characters[1][0])
	}
	if result.Characters[2][0] != 3 { // C
		t.Errorf("Row 2 first char = %d, want 3 (C)", result.Characters[2][0])
	}
}

func TestFormat_Centering(t *testing.T) {
	result := Format("HI", "note", true)

	// "HI" is 2 chars, board is 15 wide, so padding = (15-2)/2 = 6
	// H should be at position 6, but also vertically centered to row 1
	if result.Characters[1][6] != 8 { // H at center
		t.Errorf("Centered H at wrong position, got %d at [1][6]", result.Characters[1][6])
	}
	if result.Characters[1][7] != 9 { // I next to H
		t.Errorf("Centered I at wrong position, got %d at [1][7]", result.Characters[1][7])
	}
}

func TestFormat_Truncation(t *testing.T) {
	// Single long word should be truncated
	result := Format("ABCDEFGHIJKLMNOPQRST", "note", false)

	if len(result.Characters[0]) != 15 {
		t.Errorf("Row width = %d, want 15", len(result.Characters[0]))
	}
	// Last char should be O (15th letter)
	if result.Characters[0][14] != 15 {
		t.Errorf("Last char = %d, want 15 (O)", result.Characters[0][14])
	}
}

func TestFormat_FlagshipDimensions(t *testing.T) {
	result := Format("TEST", "flagship", false)

	if len(result.Characters) != 6 {
		t.Errorf("Flagship rows = %d, want 6", len(result.Characters))
	}
	if len(result.Characters[0]) != 22 {
		t.Errorf("Flagship cols = %d, want 22", len(result.Characters[0]))
	}
}

func TestFormat_EscapeCodes(t *testing.T) {
	result := Format("{red}{green}{blue}", "note", false)

	if result.Characters[0][0] != 63 { // red
		t.Errorf("Red code = %d, want 63", result.Characters[0][0])
	}
	if result.Characters[0][1] != 66 { // green
		t.Errorf("Green code = %d, want 66", result.Characters[0][1])
	}
	if result.Characters[0][2] != 67 { // blue
		t.Errorf("Blue code = %d, want 67", result.Characters[0][2])
	}
}

func TestFormat_NumericEscapeCodes(t *testing.T) {
	result := Format("{63}{64}{65}", "note", false)

	if result.Characters[0][0] != 63 {
		t.Errorf("Code 63 = %d, want 63", result.Characters[0][0])
	}
	if result.Characters[0][1] != 64 {
		t.Errorf("Code 64 = %d, want 64", result.Characters[0][1])
	}
	if result.Characters[0][2] != 65 {
		t.Errorf("Code 65 = %d, want 65", result.Characters[0][2])
	}
}

func TestFormat_HeartEscape(t *testing.T) {
	result := Format("I {<3} U", "note", false)

	if result.Characters[0][0] != 9 { // I
		t.Errorf("I = %d, want 9", result.Characters[0][0])
	}
	if result.Characters[0][2] != 62 { // heart
		t.Errorf("Heart = %d, want 62", result.Characters[0][2])
	}
	if result.Characters[0][4] != 21 { // U
		t.Errorf("U = %d, want 21", result.Characters[0][4])
	}
}

func TestFormat_AutoWrap(t *testing.T) {
	// "HELLO WORLD TEST" should wrap to multiple lines
	result := Format("HELLO WORLD TEST", "note", false)

	// First line should have "HELLO WORLD"
	if result.Characters[0][0] != 8 { // H
		t.Errorf("Row 0 first char = %d, want 8 (H)", result.Characters[0][0])
	}
	// Second line should have "TEST"
	if result.Characters[1][0] != 20 { // T
		t.Errorf("Row 1 first char = %d, want 20 (T)", result.Characters[1][0])
	}
}

func TestFormat_AutoWrapWordBoundary(t *testing.T) {
	// Should wrap at word boundaries
	result := Format("AAA BBB CCC DDD", "note", false)

	// Check that words aren't split mid-word
	// Row 0 should be "AAA BBB CCC" (11 chars)
	// Row 1 should be "DDD"
	if result.Characters[0][0] != 1 { // A
		t.Errorf("Row 0 first char = %d, want 1 (A)", result.Characters[0][0])
	}
}

func TestFormat_WarningOnTruncation(t *testing.T) {
	// Message that exceeds 45 chars (3x15) should warn
	longMsg := "THIS IS A VERY LONG MESSAGE THAT EXCEEDS THE DISPLAY"
	result := Format(longMsg, "note", false)

	if result.Warning == "" {
		t.Error("Expected warning for truncated message")
	}
}

func TestFormat_NoWarningShortMessage(t *testing.T) {
	result := Format("SHORT", "note", false)

	if result.Warning != "" {
		t.Errorf("Unexpected warning: %s", result.Warning)
	}
}

func TestGetDimensions(t *testing.T) {
	tests := []struct {
		device     string
		wantHeight int
		wantWidth  int
	}{
		{"note", 3, 15},
		{"flagship", 6, 22},
		{"", 6, 22},        // default
		{"unknown", 6, 22}, // default
	}

	for _, tt := range tests {
		h, w := getDimensions(tt.device)
		if h != tt.wantHeight || w != tt.wantWidth {
			t.Errorf("getDimensions(%q) = (%d, %d), want (%d, %d)",
				tt.device, h, w, tt.wantHeight, tt.wantWidth)
		}
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		text  string
		width int
		want  int // number of lines
	}{
		{"HELLO", 15, 1},
		{"HELLO WORLD", 15, 1},
		{"HELLO WORLD TEST", 15, 2},
		{"A B C D E F G H", 5, 3}, // "A B C", "D E F", "G H"
		{"", 15, 1},
	}

	for _, tt := range tests {
		got := wrapText(tt.text, tt.width)
		if len(got) != tt.want {
			t.Errorf("wrapText(%q, %d) = %d lines, want %d", tt.text, tt.width, len(got), tt.want)
		}
	}
}

func TestCountDisplayChars(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"HELLO", 5},
		{"HI THERE", 8},
		{"{red}", 1},            // escape code counts as 1
		{"{red}{blue}", 2},      // two escape codes
		{"A {red} B", 5},        // mixed
		{"{red}HELLO{blue}", 7}, // escape + text
		{"", 0},
	}

	for _, tt := range tests {
		got := countDisplayChars(tt.input)
		if got != tt.want {
			t.Errorf("countDisplayChars(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
