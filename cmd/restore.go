package cmd

import (
	"fmt"
	"os"

	"github.com/morheus9/k8s-backup-cli/internal/backup"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	restoreFilePath       string
	restoreNamespace      string
	restoreKubeconfigPath string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Kubernetes resources from backup",
	Long:  "Restore Kubernetes resources from a previously created backup archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		if restoreFilePath == "" {
			return fmt.Errorf("backup file path is required, use --file or -f")
		}

		// Build rest.Config to pass into restore engine.
		var (
			config *rest.Config
			err    error
		)

		if restoreKubeconfigPath != "" {
			config, err = clientcmd.BuildConfigFromFlags("", restoreKubeconfigPath)
		} else {
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			configOverrides := &clientcmd.ConfigOverrides{}
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				loadingRules,
				configOverrides,
			).ClientConfig()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building Kubernetes config: %v\n", err)
			os.Exit(1)
		}

		if err := backup.RestoreNamespace(restoreFilePath, restoreKubeconfigPath, restoreNamespace, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error restoring backup: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully restored resources from %s\n", restoreFilePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&restoreFilePath, "file", "f", "", "Path to backup archive (tar.gz) to restore from (required)")
	restoreCmd.Flags().StringVarP(&restoreNamespace, "namespace", "n", "", "Default namespace for namespaceless manifests")
	restoreCmd.Flags().StringVarP(&restoreKubeconfigPath, "kubeconfig", "k", "", "Path to kubeconfig file (default: auto-detect)")
	_ = restoreCmd.MarkFlagRequired("file")
}
