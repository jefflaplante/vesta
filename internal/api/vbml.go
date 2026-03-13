package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jeff/vesta/internal/config"
)

const vbmlURL = "https://vbml.vestaboard.com/compose"

// VBMLClient handles communication with the VBML API
type VBMLClient struct {
	httpClient *http.Client
}

// NewVBMLClient creates a new VBML API client
func NewVBMLClient() *VBMLClient {
	return &VBMLClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// VBMLRequest represents a request to the VBML compose endpoint
type VBMLRequest struct {
	Components []VBMLComponent `json:"components"`
}

// VBMLComponent represents a component in the VBML request
type VBMLComponent struct {
	Template VBMLTemplate `json:"template"`
}

// VBMLTemplate represents the template for formatting
type VBMLTemplate struct {
	Style  string `json:"style"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Justify string `json:"justify,omitempty"`
	Align   string `json:"align,omitempty"`
}

// VBMLResponse represents the response from the VBML API
type VBMLResponse struct {
	Characters [][]int `json:"characters"`
}

// Format sends a message to the VBML API for formatting
func (c *VBMLClient) Format(message string, device string, centered bool) ([][]int, error) {
	// Get device dimensions
	height, width := getDimensions(device)

	// Build request
	template := VBMLTemplate{
		Style:  message,
		Height: height,
		Width:  width,
	}

	if centered {
		template.Justify = "center"
		template.Align = "center"
	}

	req := VBMLRequest{
		Components: []VBMLComponent{
			{Template: template},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", vbmlURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VBML API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var vbmlResp VBMLResponse
	if err := json.NewDecoder(resp.Body).Decode(&vbmlResp); err != nil {
		return nil, fmt.Errorf("failed to decode VBML response: %w", err)
	}

	return vbmlResp.Characters, nil
}

// getDimensions returns the height and width for a device type
func getDimensions(device string) (int, int) {
	switch device {
	case config.DeviceFlagship:
		return 6, 22
	case config.DeviceNote:
		fallthrough
	default:
		return 3, 15
	}
}
