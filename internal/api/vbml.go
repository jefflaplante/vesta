package api

import (
	"strings"
	"unicode"

	"github.com/jeff/vesta/internal/config"
)

// Format converts a message string to a Vestaboard character array.
// Handles text, escape codes like {63}, and centering.
func Format(message string, device string, centered bool) [][]int {
	height, width := getDimensions(device)

	// Handle literal \n escape sequences
	message = strings.ReplaceAll(message, "\\n", "\n")

	// Split message into lines
	lines := strings.Split(message, "\n")

	// Limit to board height
	if len(lines) > height {
		lines = lines[:height]
	}

	// Convert each line to character codes
	var rows [][]int
	for _, line := range lines {
		row := lineToCharCodes(line, width, centered)
		rows = append(rows, row)
	}

	// Pad with empty rows if needed
	for len(rows) < height {
		rows = append(rows, make([]int, width))
	}

	// If centering vertically, shift rows
	if centered && len(lines) < height {
		emptyRows := height - len(lines)
		topPad := emptyRows / 2
		if topPad > 0 {
			newRows := make([][]int, height)
			for i := 0; i < height; i++ {
				if i < topPad || i >= topPad+len(lines) {
					newRows[i] = make([]int, width)
				} else {
					newRows[i] = rows[i-topPad]
				}
			}
			rows = newRows
		}
	}

	return rows
}

// lineToCharCodes converts a single line to character codes
func lineToCharCodes(line string, width int, centered bool) []int {
	codes := parseLineToCharCodes(line)

	// Truncate if too long
	if len(codes) > width {
		codes = codes[:width]
	}

	// Center or left-align
	if centered && len(codes) < width {
		padding := (width - len(codes)) / 2
		result := make([]int, width)
		for i, code := range codes {
			result[padding+i] = code
		}
		return result
	}

	// Pad to width
	for len(codes) < width {
		codes = append(codes, 0)
	}

	return codes
}

// parseLineToCharCodes parses a line with escape codes into character codes
func parseLineToCharCodes(line string) []int {
	var codes []int
	runes := []rune(line)
	i := 0

	for i < len(runes) {
		// Check for escape code {N} or {name}
		if runes[i] == '{' {
			end := -1
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '}' {
					end = j
					break
				}
			}
			if end > i+1 {
				content := string(runes[i+1 : end])
				if code, ok := parseEscapeCode(content); ok {
					codes = append(codes, code)
					i = end + 1
					continue
				}
			}
		}

		// Regular character
		code := charToCode(runes[i])
		codes = append(codes, code)
		i++
	}

	return codes
}

// parseEscapeCode parses the content inside braces and returns the character code
func parseEscapeCode(content string) (int, bool) {
	// Try numeric first
	var num int
	if _, err := parseNumber(content); err == nil {
		num, _ = parseNumber(content)
		if num >= 0 && num <= 71 {
			return num, true
		}
	}

	// Named codes
	switch strings.ToLower(content) {
	case "red":
		return 63, true
	case "orange":
		return 64, true
	case "yellow":
		return 65, true
	case "green":
		return 66, true
	case "blue":
		return 67, true
	case "violet":
		return 68, true
	case "white":
		return 69, true
	case "black":
		return 70, true
	case "filled":
		return 71, true
	case "deg":
		return 62, true
	case "<3":
		return 62, true
	}

	return 0, false
}

func parseNumber(s string) (int, error) {
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, nil
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}

// charToCode converts a rune to a Vestaboard character code
func charToCode(r rune) int {
	r = unicode.ToUpper(r)

	switch {
	case r >= 'A' && r <= 'Z':
		return int(r - 'A' + 1)
	case r >= '1' && r <= '9':
		return int(r - '1' + 27) // 1-9 map to codes 27-35
	case r == '0':
		return 36 // 0 maps to code 36
	case r == ' ':
		return 0
	case r == '!':
		return 37
	case r == '@':
		return 38
	case r == '#':
		return 39
	case r == '$':
		return 40
	case r == '(':
		return 41
	case r == ')':
		return 42
	case r == '-':
		return 44
	case r == '+':
		return 45
	case r == '&':
		return 46
	case r == '=':
		return 47
	case r == ';':
		return 48
	case r == ':':
		return 49
	case r == '\'':
		return 52
	case r == '"':
		return 53
	case r == '%':
		return 54
	case r == ',':
		return 55
	case r == '.':
		return 56
	case r == '/':
		return 59
	case r == '?':
		return 60
	case r == '°':
		return 62
	default:
		return 0 // Unknown characters become blank
	}
}

// getDimensions returns the height and width for a device type
func getDimensions(device string) (int, int) {
	switch device {
	case config.DeviceFlagship:
		return 6, 22
	case config.DeviceNote:
		return 3, 15
	default:
		return 6, 22
	}
}
