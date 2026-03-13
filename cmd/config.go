package cmd

import (
	"fmt"

	"github.com/jeff/vesta/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage vesta configuration",
	Long:  `View and modify vesta configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings (token is masked).`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Token:  %s\n", cfg.MaskedToken())
		fmt.Printf("Device: %s\n", cfg.Device)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  token  - Your Vestaboard API token
  device - Device type (note or flagship)

Examples:
  vesta config set token abc123
  vesta config set device flagship`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		switch key {
		case "token":
			cfg.Token = value
		case "device":
			if !config.ValidDevice(value) {
				return fmt.Errorf("invalid device type '%s'. Use 'note' or 'flagship'", value)
			}
			cfg.Device = value
		default:
			return fmt.Errorf("unknown config key '%s'. Use 'token' or 'device'", key)
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
