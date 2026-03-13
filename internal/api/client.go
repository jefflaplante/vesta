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

	if resp.StatusCode == 429 {
		return fmt.Errorf("rate limited. Vestaboard allows ~1 message per 15 seconds")
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("authentication failed. Check your API token")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// ReadResponse represents the response from reading the board
type ReadResponse struct {
	CurrentMessage struct {
		Layout [][]int `json:"layout"`
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var readResp ReadResponse
	if err := json.NewDecoder(resp.Body).Decode(&readResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return readResp.CurrentMessage.Layout, nil
}

// CharToDisplay converts a character code back to a displayable character
func CharToDisplay(code int) string {
	switch {
	case code == 0:
		return " " // blank
	case code >= 1 && code <= 26:
		return string(rune('A' + code - 1))
	case code >= 27 && code <= 36:
		return string(rune('0' + code - 27))
	case code == 37:
		return "!"
	case code == 38:
		return "@"
	case code == 39:
		return "#"
	case code == 40:
		return "$"
	case code == 41:
		return "("
	case code == 42:
		return ")"
	case code == 43:
		return "-"
	case code == 44:
		return "+"
	case code == 45:
		return "&"
	case code == 46:
		return "="
	case code == 47:
		return ";"
	case code == 48:
		return ":"
	case code == 49:
		return "'"
	case code == 50:
		return "\""
	case code == 51:
		return "%"
	case code == 52:
		return ","
	case code == 53:
		return "."
	case code == 54:
		return "/"
	case code == 55:
		return "?"
	case code == 56:
		return "°" // degree symbol display
	case code == 59:
		return "🔴" // red
	case code == 60:
		return "🟠" // orange
	case code == 61:
		return "🟡" // yellow
	case code == 62:
		return "♥" // heart/degree (context dependent)
	case code == 63:
		return "🔴" // red
	case code == 64:
		return "🟠" // orange
	case code == 65:
		return "🟡" // yellow
	case code == 66:
		return "🟢" // green
	case code == 67:
		return "🔵" // blue
	case code == 68:
		return "🟣" // violet
	case code == 69:
		return "⬜" // white
	case code == 70:
		return "⬛" // black
	case code == 71:
		return "█" // filled
	default:
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
