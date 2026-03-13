package cmd

import (
	"fmt"

	"github.com/jeff/vesta/internal/api"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read the current board state",
	Long:  `Read and display the current message on your Vestaboard.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create client from config
		client, err := api.NewClientFromConfig(cfg)
		if err != nil {
			return err
		}

		// Read from board
		layout, err := client.Read()
		if err != nil {
			return err
		}

		// Display the board
		fmt.Println(api.DisplayBoard(layout))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
