package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token      string `yaml:"token"`
	Device     string `yaml:"device"`
	LocalURL   string `yaml:"local_url"`   // e.g., "192.168.1.100:7000"
	LocalToken string `yaml:"local_token"` // Local API key
	APIMode    string `yaml:"api_mode"`    // "cloud" (default) or "local"
}

const (
	DeviceNote     = "note"
	DeviceFlagship = "flagship"
	APIModeCloud   = "cloud"
	APIModeLocal   = "local"
)

var configDir string
var configFile string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	configDir = filepath.Join(home, ".config", "vesta")
	configFile = filepath.Join(configDir, "config.yaml")
}

// Load reads config from file and environment, with proper precedence.
// Priority: env var > config file > defaults
func Load() (*Config, error) {
	cfg := &Config{
		Device:  DeviceNote,  // default
		APIMode: APIModeCloud, // default
	}

	// Try to load from config file
	if data, err := os.ReadFile(configFile); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Environment variables override config file
	if envToken := os.Getenv("VESTABOARD_API_KEY"); envToken != "" {
		cfg.Token = envToken
		os.Unsetenv("VESTABOARD_API_KEY")
	}
	if envLocalToken := os.Getenv("VESTABOARD_LOCAL_API_KEY"); envLocalToken != "" {
		cfg.LocalToken = envLocalToken
		os.Unsetenv("VESTABOARD_LOCAL_API_KEY")
	}
	if envLocalURL := os.Getenv("VESTABOARD_LOCAL_URL"); envLocalURL != "" {
		cfg.LocalURL = envLocalURL
	}
	if envAPIMode := os.Getenv("VESTABOARD_API_MODE"); envAPIMode != "" {
		cfg.APIMode = envAPIMode
	}

	// Default APIMode if empty
	if cfg.APIMode == "" {
		cfg.APIMode = APIModeCloud
	}

	return cfg, nil
}

// Save writes the config to the config file
func (c *Config) Save() error {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetToken returns the token, with a helpful error if not configured
func (c *Config) GetToken() (string, error) {
	if c.Token == "" {
		return "", fmt.Errorf("no API token configured. Run 'vesta config set token <key>' or set VESTABOARD_API_KEY")
	}
	return c.Token, nil
}

// MaskedToken returns the token with all but the last 4 characters masked
func (c *Config) MaskedToken() string {
	if c.Token == "" {
		return "(not set)"
	}
	if len(c.Token) <= 4 {
		return "****"
	}
	return "****" + c.Token[len(c.Token)-4:]
}

// ValidDevice checks if the device type is valid
func ValidDevice(device string) bool {
	return device == DeviceNote || device == DeviceFlagship
}

// ValidAPIMode checks if the API mode is valid
func ValidAPIMode(mode string) bool {
	return mode == APIModeCloud || mode == APIModeLocal
}

// GetLocalToken returns the local API token, with a helpful error if not configured
func (c *Config) GetLocalToken() (string, error) {
	if c.LocalToken == "" {
		return "", fmt.Errorf("no local API token configured. Run 'vesta config set local-token' or set VESTABOARD_LOCAL_API_KEY")
	}
	return c.LocalToken, nil
}

// GetLocalURL returns the local API URL, with a helpful error if not configured
func (c *Config) GetLocalURL() (string, error) {
	if c.LocalURL == "" {
		return "", fmt.Errorf("no local URL configured. Run 'vesta config set local-url <ip:port>' or set VESTABOARD_LOCAL_URL")
	}
	return c.LocalURL, nil
}

// MaskedLocalToken returns the local token with all but the last 4 characters masked
func (c *Config) MaskedLocalToken() string {
	if c.LocalToken == "" {
		return "(not set)"
	}
	if len(c.LocalToken) <= 4 {
		return "****"
	}
	return "****" + c.LocalToken[len(c.LocalToken)-4:]
}

// IsLocalMode returns true if the API mode is set to local
func (c *Config) IsLocalMode() bool {
	return c.APIMode == APIModeLocal
}
