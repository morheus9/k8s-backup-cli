package list

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/morheus9/k8s-backup-cli/internal/k8s"
)

// ListResources lists all resources in the specified namespace
func ListResources(namespace string, kubeconfigPath string) error {
	var client *k8s.Client
	var err error

	if kubeconfigPath != "" {
		client, err = k8s.NewClient(kubeconfigPath)
	} else {
		client, err = k8s.NewClientFromDefault()
	}
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	ctx := context.Background()
	resources, err := client.FetchResources(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to fetch resources: %w", err)
	}

	if len(resources) == 0 {
		fmt.Printf("No resources found in namespace '%s'\n", namespace)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintf(w, "KIND\tNAME\tNAMESPACE\tAPI VERSION\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintf(w, "----\t----\t---------\t-----------\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	for _, resource := range resources {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			resource.Kind,
			resource.Name,
			resource.Namespace,
			resource.APIVersion,
		); err != nil {
			return fmt.Errorf("failed to write resource: %w", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}
	fmt.Printf("\nTotal: %d resources\n", len(resources))

	return nil
}
