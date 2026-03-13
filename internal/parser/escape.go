package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeff/vesta/internal/config"
)

// Character codes for Vestaboard
const (
	CodeDegree = 62 // Flagship only
	CodeHeart  = 62 // Note only (same code, different meaning)
	CodeRed    = 63
	CodeOrange = 64
	CodeYellow = 65
	CodeGreen  = 66
	CodeBlue   = 67
	CodeViolet = 68
	CodeWhite  = 69
	CodeBlack  = 70
	CodeFilled = 71
)

// Named escape codes (without the braces)
var namedCodes = map[string]int{
	"red":    CodeRed,
	"orange": CodeOrange,
	"yellow": CodeYellow,
	"green":  CodeGreen,
	"blue":   CodeBlue,
	"violet": CodeViolet,
	"white":  CodeWhite,
	"black":  CodeBlack,
	"filled": CodeFilled,
	"deg":    CodeDegree,
	"<3":     CodeHeart,
}

// DeviceOnlySymbols maps codes to their required device type
var deviceOnlySymbols = map[string]string{
	"deg": config.DeviceFlagship,
	"<3":  config.DeviceNote,
}

// Warning represents a parser warning
type Warning struct {
	Symbol  string
	Message string
}

// ParseResult contains the parsed message and any warnings
type ParseResult struct {
	Message  string
	Warnings []Warning
}

// Parse converts user escape syntax to VBML format.
// Supports: {color}, {deg}, {<3}, bare <3, and {N} for raw codes.
// Returns the converted message and any warnings about device compatibility.
func Parse(input string, device string) ParseResult {
	result := ParseResult{
		Message:  input,
		Warnings: []Warning{},
	}

	// First, convert bare <3 (not in braces) to {62}
	// But don't convert if it's already inside braces
	bareHeartRe := regexp.MustCompile(`(?:^|[^{])<3(?:[^}]|$)`)
	result.Message = bareHeartRe.ReplaceAllStringFunc(result.Message, func(match string) string {
		// Preserve characters before and after <3
		prefix := ""
		suffix := ""
		if len(match) > 2 && match[0] != '<' {
			prefix = string(match[0])
		}
		if len(match) > 2 && match[len(match)-1] != '3' {
			suffix = string(match[len(match)-1])
		}
		return prefix + "{62}" + suffix
	})

	// Now process all {xxx} codes
	braceRe := regexp.MustCompile(`\{([^}]+)\}`)
	result.Message = braceRe.ReplaceAllStringFunc(result.Message, func(match string) string {
		// Extract content between braces
		content := match[1 : len(match)-1]
		content = strings.ToLower(content)

		// Check for named codes
		if code, ok := namedCodes[content]; ok {
			// Check device compatibility
			if requiredDevice, hasRestriction := deviceOnlySymbols[content]; hasRestriction {
				if device != requiredDevice {
					var friendlyName string
					if content == "deg" {
						friendlyName = "degree symbol"
					} else {
						friendlyName = "heart"
					}
					result.Warnings = append(result.Warnings, Warning{
						Symbol:  content,
						Message: fmt.Sprintf("{%s} (%s) is only supported on Vestaboard %s", content, friendlyName, strings.Title(requiredDevice)),
					})
				}
			}
			return fmt.Sprintf("{%d}", code)
		}

		// Check for numeric codes
		if num, err := strconv.Atoi(content); err == nil {
			if num >= 0 && num <= 71 {
				return fmt.Sprintf("{%d}", num)
			}
			result.Warnings = append(result.Warnings, Warning{
				Symbol:  content,
				Message: fmt.Sprintf("character code %d is out of range (0-71)", num),
			})
		}

		// Unknown code - leave as-is (VBML will handle it)
		return match
	})

	return result
}

// FormatWarnings returns warnings as a string for display
func FormatWarnings(warnings []Warning) string {
	if len(warnings) == 0 {
		return ""
	}
	var lines []string
	for _, w := range warnings {
		lines = append(lines, fmt.Sprintf("Warning: %s", w.Message))
	}
	return strings.Join(lines, "\n")
}
