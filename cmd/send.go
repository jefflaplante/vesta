package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jeff/vesta/internal/api"
	"github.com/jeff/vesta/internal/parser"
	"github.com/spf13/cobra"
)

var centerFlag bool
var dryRunFlag bool

var sendCmd = &cobra.Command{
	Use:   "send [message]",
	Short: "Send a message to your Vestaboard",
	Long: `Send a formatted message to your Vestaboard.

Supports escape codes for colors and symbols:
  {red}, {orange}, {yellow}, {green}, {blue}, {violet}, {white}, {black}, {filled}
  {deg} (Flagship only), {<3} or <3 (Note only)
  {0}-{71} for raw character codes

Use "-" to read from stdin for scripting:
  echo "Hello" | vesta send -

Examples:
  vesta send "Hello World"
  vesta send -c "Centered Message"
  vesta send "I <3 Go"
  vesta send "{red}{green}{blue}"
  echo "Piped input" | vesta send -`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]

		// Read from stdin if message is "-"
		if message == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			message = strings.TrimSuffix(string(data), "\n")
		}

		// Get token
		token, err := cfg.GetToken()
		if err != nil {
			return err
		}

		// Parse escape codes
		parseResult := parser.Parse(message, cfg.Device)

		// Print any warnings
		if warnings := parser.FormatWarnings(parseResult.Warnings); warnings != "" {
			fmt.Println(warnings)
		}

		// Format message to character array (with auto-wrap)
		formatResult := api.Format(parseResult.Message, cfg.Device, centerFlag)

		// Print format warning if any
		if formatResult.Warning != "" {
			fmt.Printf("Warning: %s\n", formatResult.Warning)
		}

		// Dry run - show what would be sent
		if dryRunFlag {
			fmt.Println("Character array:")
			for i, row := range formatResult.Characters {
				fmt.Printf("Row %d: %v\n", i, row)
			}
			fmt.Println("\nPreview:")
			fmt.Println(api.DisplayBoard(formatResult.Characters))
			return nil
		}

		// Send to board
		client := api.NewClient(token)
		if err := client.Send(formatResult.Characters); err != nil {
			if apiErr, ok := err.(*api.APIError); ok && VerboseFlag {
				fmt.Fprintln(cmd.ErrOrStderr(), apiErr.VerboseMessage())
			}
			return err
		}

		fmt.Println("Message sent successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().BoolVarP(&centerFlag, "center", "c", false, "center the message")
	sendCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "show character array without sending")
}
