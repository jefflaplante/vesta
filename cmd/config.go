package cmd

import (
	"fmt"
	"os"

	"github.com/jeff/vesta/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage vesta configuration",
	Long:  `View and modify vesta configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings (tokens are masked).`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Token:       %s\n", cfg.MaskedToken())
		fmt.Printf("Device:      %s\n", cfg.Device)
		fmt.Printf("API Mode:    %s\n", cfg.APIMode)
		fmt.Printf("Local URL:   %s\n", valueOrNotSet(cfg.LocalURL))
		fmt.Printf("Local Token: %s\n", cfg.MaskedLocalToken())
		return nil
	},
}

func valueOrNotSet(v string) string {
	if v == "" {
		return "(not set)"
	}
	return v
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  token       - Your Vestaboard cloud API token
  device      - Device type (note or flagship)
  api-mode    - API mode (cloud or local)
  local-url   - Local API URL (ip:port, e.g., 192.168.1.100:7000)
  local-token - Local API token

For security, omit the value for token/local-token to enter it interactively:
  vesta config set token
  vesta config set local-token

Examples:
  vesta config set token abc123
  vesta config set token            # prompts securely
  vesta config set device flagship
  vesta config set api-mode local
  vesta config set local-url 192.168.1.100:7000
  vesta config set local-token      # prompts securely`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		switch key {
		case "token":
			var value string
			if len(args) < 2 {
				fmt.Print("Enter API token: ")
				bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read token: %w", err)
				}
				fmt.Println()
				value = string(bytePassword)
			} else {
				value = args[1]
			}
			cfg.Token = value
		case "device":
			if len(args) < 2 {
				return fmt.Errorf("device requires a value: 'note' or 'flagship'")
			}
			value := args[1]
			if !config.ValidDevice(value) {
				return fmt.Errorf("invalid device type '%s'. Use 'note' or 'flagship'", value)
			}
			cfg.Device = value
		case "api-mode":
			if len(args) < 2 {
				return fmt.Errorf("api-mode requires a value: 'cloud' or 'local'")
			}
			value := args[1]
			if !config.ValidAPIMode(value) {
				return fmt.Errorf("invalid API mode '%s'. Use 'cloud' or 'local'", value)
			}
			cfg.APIMode = value
		case "local-url":
			if len(args) < 2 {
				return fmt.Errorf("local-url requires a value (e.g., 192.168.1.100:7000)")
			}
			cfg.LocalURL = args[1]
		case "local-token":
			var value string
			if len(args) < 2 {
				fmt.Print("Enter local API token: ")
				bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return fmt.Errorf("failed to read token: %w", err)
				}
				fmt.Println()
				value = string(bytePassword)
			} else {
				value = args[1]
			}
			cfg.LocalToken = value
		default:
			return fmt.Errorf("unknown config key '%s'. Use 'token', 'device', 'api-mode', 'local-url', or 'local-token'", key)
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("Configuration updated: %s\n", key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}
