package parser

import (
	"testing"
)

func TestParse_NoEscapes(t *testing.T) {
	result := Parse("Hello World", "note")

	if result.Message != "Hello World" {
		t.Errorf("Message = %q, want \"Hello World\"", result.Message)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("Warnings = %v, want none", result.Warnings)
	}
}

func TestParse_ColorEscapes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"{red}", "{63}"},
		{"{orange}", "{64}"},
		{"{yellow}", "{65}"},
		{"{green}", "{66}"},
		{"{blue}", "{67}"},
		{"{violet}", "{68}"},
		{"{white}", "{69}"},
		{"{black}", "{70}"},
		{"{filled}", "{71}"},
	}

	for _, tt := range tests {
		result := Parse(tt.input, "note")
		if result.Message != tt.want {
			t.Errorf("Parse(%q) = %q, want %q", tt.input, result.Message, tt.want)
		}
	}
}

func TestParse_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"{RED}", "{63}"},
		{"{Red}", "{63}"},
		{"{rEd}", "{63}"},
	}

	for _, tt := range tests {
		result := Parse(tt.input, "note")
		if result.Message != tt.want {
			t.Errorf("Parse(%q) = %q, want %q", tt.input, result.Message, tt.want)
		}
	}
}

func TestParse_NumericEscapes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"{0}", "{0}"},
		{"{63}", "{63}"},
		{"{71}", "{71}"},
	}

	for _, tt := range tests {
		result := Parse(tt.input, "note")
		if result.Message != tt.want {
			t.Errorf("Parse(%q) = %q, want %q", tt.input, result.Message, tt.want)
		}
	}
}

func TestParse_HeartEscape(t *testing.T) {
	result := Parse("{<3}", "note")
	if result.Message != "{62}" {
		t.Errorf("Parse(\"{<3}\") = %q, want \"{62}\"", result.Message)
	}
}

func TestParse_DegreeEscape(t *testing.T) {
	result := Parse("{deg}", "flagship")
	if result.Message != "{62}" {
		t.Errorf("Parse(\"{deg}\") = %q, want \"{62}\"", result.Message)
	}
}

func TestParse_DegreeWarningOnNote(t *testing.T) {
	result := Parse("{deg}", "note")

	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
		return
	}
	if result.Warnings[0].Symbol != "deg" {
		t.Errorf("Warning symbol = %q, want \"deg\"", result.Warnings[0].Symbol)
	}
}

func TestParse_HeartWarningOnFlagship(t *testing.T) {
	result := Parse("{<3}", "flagship")

	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
		return
	}
	if result.Warnings[0].Symbol != "<3" {
		t.Errorf("Warning symbol = %q, want \"<3\"", result.Warnings[0].Symbol)
	}
}

func TestParse_MixedContent(t *testing.T) {
	result := Parse("Hello {red} World {blue}", "note")

	want := "Hello {63} World {67}"
	if result.Message != want {
		t.Errorf("Parse mixed = %q, want %q", result.Message, want)
	}
}

func TestParse_MultipleEscapes(t *testing.T) {
	result := Parse("{red}{green}{blue}", "note")

	want := "{63}{66}{67}"
	if result.Message != want {
		t.Errorf("Parse multiple = %q, want %q", result.Message, want)
	}
}

func TestParse_UnknownEscapePassthrough(t *testing.T) {
	result := Parse("{unknown}", "note")

	// Unknown escapes should pass through unchanged
	if result.Message != "{unknown}" {
		t.Errorf("Parse unknown = %q, want \"{unknown}\"", result.Message)
	}
}

func TestParse_OutOfRangeNumeric(t *testing.T) {
	result := Parse("{72}", "note")

	// Out of range should generate warning and pass through
	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning for out of range, got %d", len(result.Warnings))
	}
}

func TestFormatWarnings_Empty(t *testing.T) {
	result := FormatWarnings(nil)
	if result != "" {
		t.Errorf("FormatWarnings(nil) = %q, want empty", result)
	}

	result = FormatWarnings([]Warning{})
	if result != "" {
		t.Errorf("FormatWarnings([]) = %q, want empty", result)
	}
}

func TestFormatWarnings_Single(t *testing.T) {
	warnings := []Warning{{Symbol: "deg", Message: "test message"}}
	result := FormatWarnings(warnings)

	if result != "Warning: test message" {
		t.Errorf("FormatWarnings = %q, want \"Warning: test message\"", result)
	}
}

func TestFormatWarnings_Multiple(t *testing.T) {
	warnings := []Warning{
		{Symbol: "a", Message: "first"},
		{Symbol: "b", Message: "second"},
	}
	result := FormatWarnings(warnings)

	want := "Warning: first\nWarning: second"
	if result != want {
		t.Errorf("FormatWarnings = %q, want %q", result, want)
	}
}
