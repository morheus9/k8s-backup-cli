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
	Use:   "backup [namespace]",
	Short: "Create a backup of Kubernetes resources",
	Long:  "Create a tar.gz archive with Kubernetes manifests from the specified namespace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := backupNamespace
		if len(args) > 0 {
			ns = args[0]
		}

		if ns == "" {
			return fmt.Errorf("namespace is required. Use --namespace flag or provide as argument")
		}

		archivePath, err := backup.BackupNamespace(ns, backupKubeconfigPath)
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
	backupCmd.Flags().StringVarP(&backupNamespace, "namespace", "n", "", "Kubernetes namespace to backup")
	backupCmd.Flags().StringVarP(&backupKubeconfigPath, "kubeconfig", "k", "", "Path to kubeconfig file (default: auto-detect)")
}
