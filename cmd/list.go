package cmd

import (
	"fmt"
	"os"

	"github.com/morheus9/k8s-backup-cli/internal/list"
	"github.com/spf13/cobra"
)

var (
	namespace      string
	kubeconfigPath string
)

var listCmd = &cobra.Command{
	Use:   "list [namespace]",
	Short: "List Kubernetes resources in namespace",
	Long:  "List all Kubernetes resources in the specified namespace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := namespace
		if len(args) > 0 {
			ns = args[0]
		}

		if ns == "" {
			return fmt.Errorf("namespace is required. Use --namespace flag or provide as argument")
		}

		if err := list.ListResources(ns, kubeconfigPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	listCmd.Flags().StringVarP(&kubeconfigPath, "kubeconfig", "k", "", "Path to kubeconfig file (default: use default kubeconfig)")
}
