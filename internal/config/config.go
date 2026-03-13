package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token  string `yaml:"token"`
	Device string `yaml:"device"`
}

const (
	DeviceNote     = "note"
	DeviceFlagship = "flagship"
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
		Device: DeviceNote, // default
	}

	// Try to load from config file
	if data, err := os.ReadFile(configFile); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Environment variable overrides config file
	if envToken := os.Getenv("VESTABOARD_API_KEY"); envToken != "" {
		cfg.Token = envToken
		os.Unsetenv("VESTABOARD_API_KEY")
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
