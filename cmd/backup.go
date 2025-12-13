package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of Kubernetes resources",
	Long:  "Create a backup of Kubernetes resources from the specified namespace or cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Backup command executed")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Add flags here if needed
	// backupCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace to backup")
}
