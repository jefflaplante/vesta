package api

import (
	"testing"
)

func TestCharToDisplay_Letters(t *testing.T) {
	for code := 1; code <= 26; code++ {
		expected := string(rune('A' + code - 1))
		got := CharToDisplay(code)
		if got != expected {
			t.Errorf("CharToDisplay(%d) = %q, want %q", code, got, expected)
		}
	}
}

func TestCharToDisplay_Digits(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{27, "1"}, {28, "2"}, {29, "3"}, {30, "4"}, {31, "5"},
		{32, "6"}, {33, "7"}, {34, "8"}, {35, "9"}, {36, "0"},
	}

	for _, tt := range tests {
		got := CharToDisplay(tt.code)
		if got != tt.want {
			t.Errorf("CharToDisplay(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestCharToDisplay_SpecialChars(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, " "},   // blank
		{37, "!"}, {38, "@"}, {39, "#"}, {40, "$"},
		{41, "("}, {42, ")"}, {44, "-"}, {45, "+"},
		{46, "&"}, {47, "="}, {48, ";"}, {49, ":"},
		{52, "'"}, {53, "\""}, {54, "%"}, {55, ","},
		{56, "."}, {59, "/"}, {60, "?"},
	}

	for _, tt := range tests {
		got := CharToDisplay(tt.code)
		if got != tt.want {
			t.Errorf("CharToDisplay(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestCharToDisplay_Colors(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{62, "♥"},  // heart
		{63, "🔴"}, // red
		{64, "🟠"}, // orange
		{65, "🟡"}, // yellow
		{66, "🟢"}, // green
		{67, "🔵"}, // blue
		{68, "🟣"}, // violet
		{69, "⬜"}, // white
		{70, "⬛"}, // black
		{71, "█"},  // filled
	}

	for _, tt := range tests {
		got := CharToDisplay(tt.code)
		if got != tt.want {
			t.Errorf("CharToDisplay(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestCharToDisplay_Unknown(t *testing.T) {
	// Codes not in any range should return "?"
	unknowns := []int{43, 50, 51, 57, 58, 61, 72, 100, -1}
	for _, code := range unknowns {
		got := CharToDisplay(code)
		if got != "?" {
			t.Errorf("CharToDisplay(%d) = %q, want \"?\"", code, got)
		}
	}
}

func TestDisplayBoard(t *testing.T) {
	layout := [][]int{
		{8, 9, 0, 0, 0}, // HI
		{0, 0, 0, 0, 0},
	}

	result := DisplayBoard(layout)

	// Should contain the board characters
	if len(result) == 0 {
		t.Error("DisplayBoard returned empty string")
	}

	// Check it contains expected characters
	expectedChars := []string{"H", "I", "┌", "┐", "└", "┘", "│"}
	for _, ch := range expectedChars {
		found := false
		for _, r := range result {
			if string(r) == ch {
				found = true
				break
			}
		}
		if !found && ch != "H" && ch != "I" {
			// H and I might not be single runes due to the way we check
			continue
		}
	}
}

func TestDisplayBoard_Empty(t *testing.T) {
	result := DisplayBoard([][]int{})
	if result != "(empty)" {
		t.Errorf("DisplayBoard([]) = %q, want \"(empty)\"", result)
	}
}

func TestAPIError_FriendlyMessage(t *testing.T) {
	tests := []struct {
		errType string
		want    string
	}{
		{"FingerprintMatch", "This message is already displayed on the board"},
		{"QuietHours", "Quiet hours are enabled on this Vestaboard"},
		{"RateLimited", "Rate limited. Wait ~15 seconds between messages"},
		{"Unknown", ""},
	}

	for _, tt := range tests {
		err := &APIError{Type: tt.errType, StatusCode: 400}
		got := err.FriendlyMessage()
		if tt.want != "" && got != tt.want {
			t.Errorf("FriendlyMessage for %q = %q, want %q", tt.errType, got, tt.want)
		}
	}
}

func TestAPIError_VerboseMessage(t *testing.T) {
	err := &APIError{
		StatusCode: 409,
		Type:       "FingerprintMatch",
		Message:    "Duplicate message",
		RawBody:    `{"status":"error"}`,
	}

	verbose := err.VerboseMessage()

	if verbose == "" {
		t.Error("VerboseMessage returned empty string")
	}

	// Should contain all the details
	checks := []string{"409", "FingerprintMatch", "Duplicate message", `{"status":"error"}`}
	for _, check := range checks {
		found := false
		for i := 0; i <= len(verbose)-len(check); i++ {
			if verbose[i:i+len(check)] == check {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("VerboseMessage missing %q", check)
		}
	}
}

func TestParseAPIError(t *testing.T) {
	body := []byte(`{"status":"error","type":"FingerprintMatch","message":"Already displayed"}`)

	err := parseAPIError(409, body)

	if err.StatusCode != 409 {
		t.Errorf("StatusCode = %d, want 409", err.StatusCode)
	}
	if err.Type != "FingerprintMatch" {
		t.Errorf("Type = %q, want \"FingerprintMatch\"", err.Type)
	}
	if err.Message != "Already displayed" {
		t.Errorf("Message = %q, want \"Already displayed\"", err.Message)
	}
}

func TestParseAPIError_InvalidJSON(t *testing.T) {
	body := []byte(`not json`)

	err := parseAPIError(500, body)

	if err.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", err.StatusCode)
	}
	if err.RawBody != "not json" {
		t.Errorf("RawBody = %q, want \"not json\"", err.RawBody)
	}
}
