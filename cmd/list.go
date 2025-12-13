package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available backups",
	Long:  "List all available backups stored locally or remotely",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation here
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Add flags here if needed
}
