package cmd

import (
	"fmt"
	"os"

	"github.com/morheus9/k8s-backup-cli/internal/backup"
	"github.com/spf13/cobra"
)

var (
	backupNamespace      string
	backupKubeconfigPath string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of Kubernetes resources",
	Long:  "Create a tar.gz archive with Kubernetes manifests from the specified namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if backupNamespace == "" {
			return fmt.Errorf("namespace is required, use --namespace or -n")
		}

		archivePath, err := backup.BackupNamespace(backupNamespace, backupKubeconfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating backup: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Backup created at: %s\n", archivePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&backupNamespace, "namespace", "n", "", "Kubernetes namespace to backup (required)")
	backupCmd.Flags().StringVarP(&backupKubeconfigPath, "kubeconfig", "k", "", "Path to kubeconfig file (default: auto-detect)")
	_ = backupCmd.MarkFlagRequired("namespace")
}
