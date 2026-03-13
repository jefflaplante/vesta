package cmd

import (
	"fmt"
	"os"

	"github.com/jeff/vesta/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	deviceFlag string
)

var rootCmd = &cobra.Command{
	Use:   "vesta",
	Short: "Vestaboard CLI - send formatted messages to your Vestaboard",
	Long: `vesta is a command-line tool for sending messages to Vestaboard Note and Flagship devices.

Supports escape codes for colors and symbols:
  {red}, {orange}, {yellow}, {green}, {blue}, {violet}, {white}, {black}, {filled}
  {deg} (Flagship only), {<3} or <3 (Note only)
  {0}-{71} for raw character codes

Examples:
  vesta send "Hello World"
  vesta send -c "Centered"
  vesta send "I <3 Go"
  vesta send "{red}{green}{blue} Colors"`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Override device from flag
		if deviceFlag != "" {
			if !config.ValidDevice(deviceFlag) {
				return fmt.Errorf("invalid device type '%s'. Use 'note' or 'flagship'", deviceFlag)
			}
			cfg.Device = deviceFlag
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&deviceFlag, "device", "", "device type (note or flagship)")
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}
