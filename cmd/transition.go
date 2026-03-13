package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var transitionCmd = &cobra.Command{
	Use:   "transition",
	Short: "Manage board transition settings",
	Long:  `Get or set the transition effect used when updating the board.`,
}

var transitionGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current transition settings",
	Long:  `Display the current transition effect and speed settings.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Note: The Vestaboard Read/Write API doesn't expose transition settings.
		// This would require the Subscription API which uses a different auth mechanism.
		fmt.Println("Transition settings are managed through the Vestaboard app.")
		fmt.Println("The Read/Write API does not support reading transition settings.")
		return nil
	},
}

var transitionSpeed string

var transitionSetCmd = &cobra.Command{
	Use:   "set [type]",
	Short: "Set transition effect",
	Long: `Set the transition effect used when updating the board.

Available types: classic, wave, drift, curtain
Available speeds: gentle, fast

Note: Transition settings require the Subscription API.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		transitionType := args[0]

		// Validate type
		validTypes := map[string]bool{
			"classic": true,
			"wave":    true,
			"drift":   true,
			"curtain": true,
		}

		if !validTypes[transitionType] {
			return fmt.Errorf("invalid transition type '%s'. Use: classic, wave, drift, or curtain", transitionType)
		}

		// Validate speed if provided
		if transitionSpeed != "" {
			validSpeeds := map[string]bool{
				"gentle": true,
				"fast":   true,
			}
			if !validSpeeds[transitionSpeed] {
				return fmt.Errorf("invalid speed '%s'. Use: gentle or fast", transitionSpeed)
			}
		}

		// Note: The Vestaboard Read/Write API doesn't support setting transitions.
		fmt.Println("Transition settings are managed through the Vestaboard app.")
		fmt.Println("The Read/Write API does not support setting transition settings.")
		fmt.Printf("Requested: type=%s", transitionType)
		if transitionSpeed != "" {
			fmt.Printf(", speed=%s", transitionSpeed)
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(transitionCmd)
	transitionCmd.AddCommand(transitionGetCmd)
	transitionCmd.AddCommand(transitionSetCmd)
	transitionSetCmd.Flags().StringVar(&transitionSpeed, "speed", "", "transition speed (gentle or fast)")
}
