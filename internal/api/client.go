package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	cloudBaseURL = "https://rw.vestaboard.com/"
)

// APIError represents an error from the Vestaboard API
type APIError struct {
	StatusCode int
	Type       string
	Message    string
	RawBody    string
}

func (e *APIError) Error() string {
	return e.FriendlyMessage()
}

// FriendlyMessage returns a user-friendly error message
func (e *APIError) FriendlyMessage() string {
	switch e.Type {
	case "FingerprintMatch":
		return "This message is already displayed on the board"
	case "QuietHours":
		return "Quiet hours are enabled on this Vestaboard"
	case "RateLimited":
		return "Rate limited. Wait ~15 seconds between messages"
	default:
		if e.Message != "" {
			return e.Message
		}
		return fmt.Sprintf("API error (status %d)", e.StatusCode)
	}
}

// VerboseMessage returns detailed error information
func (e *APIError) VerboseMessage() string {
	return fmt.Sprintf("Status: %d\nType: %s\nMessage: %s\nRaw: %s",
		e.StatusCode, e.Type, e.Message, e.RawBody)
}

// parseAPIError extracts error details from API response
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		RawBody:    string(body),
	}

	var errResp struct {
		Status  string `json:"status"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &errResp) == nil {
		apiErr.Type = errResp.Type
		apiErr.Message = errResp.Message
	}

	return apiErr
}

// Client handles communication with the Vestaboard Cloud API
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new Vestaboard API client
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send sends a character array to the Vestaboard
func (c *Client) Send(characters [][]int) error {
	// Read/Write API expects the array directly
	body, err := json.Marshal(characters)
	if err != nil {
		return fmt.Errorf("failed to marshal characters: %w", err)
	}

	req, err := http.NewRequest("POST", cloudBaseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vestaboard-Read-Write-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("authentication failed. Check your API token")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return parseAPIError(resp.StatusCode, bodyBytes)
	}

	return nil
}

// ReadResponse represents the response from reading the board
type ReadResponse struct {
	CurrentMessage struct {
		Layout json.RawMessage `json:"layout"`
	} `json:"currentMessage"`
}

// Read retrieves the current board state
func (c *Client) Read() ([][]int, error) {
	req, err := http.NewRequest("GET", cloudBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Vestaboard-Read-Write-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, fmt.Errorf("authentication failed. Check your API token")
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var readResp ReadResponse
	if err := json.Unmarshal(bodyBytes, &readResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Try to parse layout as [][]int first
	var layout [][]int
	if err := json.Unmarshal(readResp.CurrentMessage.Layout, &layout); err != nil {
		// Layout might be a string-encoded array, try parsing it
		var layoutStr string
		if err := json.Unmarshal(readResp.CurrentMessage.Layout, &layoutStr); err != nil {
			return nil, fmt.Errorf("failed to parse layout: %w", err)
		}
		if err := json.Unmarshal([]byte(layoutStr), &layout); err != nil {
			return nil, fmt.Errorf("failed to parse layout string: %w", err)
		}
	}

	return layout, nil
}

// CharToDisplay converts a character code back to a displayable character
func CharToDisplay(code int) string {
	switch code {
	case 0:
		return " "
	case 37:
		return "!"
	case 38:
		return "@"
	case 39:
		return "#"
	case 40:
		return "$"
	case 41:
		return "("
	case 42:
		return ")"
	case 44:
		return "-"
	case 45:
		return "+"
	case 46:
		return "&"
	case 47:
		return "="
	case 48:
		return ";"
	case 49:
		return ":"
	case 52:
		return "'"
	case 53:
		return "\""
	case 54:
		return "%"
	case 55:
		return ","
	case 56:
		return "."
	case 59:
		return "/"
	case 60:
		return "?"
	case 62:
		return "♥" // heart/degree
	case 63:
		return "🔴" // red
	case 64:
		return "🟠" // orange
	case 65:
		return "🟡" // yellow
	case 66:
		return "🟢" // green
	case 67:
		return "🔵" // blue
	case 68:
		return "🟣" // violet
	case 69:
		return "⬜" // white
	case 70:
		return "⬛" // black
	case 71:
		return "█" // filled
	default:
		if code >= 1 && code <= 26 {
			return string(rune('A' + code - 1))
		}
		if code >= 27 && code <= 35 {
			return string(rune('1' + code - 27))
		}
		if code == 36 {
			return "0"
		}
		return "?"
	}
}

// DisplayBoard renders the board state as a string
func DisplayBoard(layout [][]int) string {
	if len(layout) == 0 {
		return "(empty)"
	}

	var sb strings.Builder
	sb.WriteString("┌")
	for i := 0; i < len(layout[0]); i++ {
		sb.WriteString("─")
	}
	sb.WriteString("┐\n")

	for _, row := range layout {
		sb.WriteString("│")
		for _, code := range row {
			sb.WriteString(CharToDisplay(code))
		}
		sb.WriteString("│\n")
	}

	sb.WriteString("└")
	for i := 0; i < len(layout[0]); i++ {
		sb.WriteString("─")
	}
	sb.WriteString("┘")

	return sb.String()
}
