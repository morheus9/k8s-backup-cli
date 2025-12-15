package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/morheus9/k8s-backup-cli/internal/k8s"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// BackupNamespace creates a tar.gz archive with Kubernetes manifests for all supported
// resources in the given namespace. The archive is created in the current working directory.
// It returns the full path to the created archive.
func BackupNamespace(namespace, kubeconfigPath string) (string, error) {
	if namespace == "" {
		return "", fmt.Errorf("namespace is required")
	}

	var (
		client *k8s.Client
		err    error
	)

	if kubeconfigPath != "" {
		client, err = k8s.NewClient(kubeconfigPath)
	} else {
		client, err = k8s.NewClientFromDefault()
	}
	if err != nil {
		return "", fmt.Errorf("create Kubernetes client: %w", err)
	}

	ctx := context.Background()
	manifests, err := client.ExportNamespaceManifests(ctx, namespace)
	if err != nil {
		return "", fmt.Errorf("export manifests: %w", err)
	}

	if len(manifests) == 0 {
		return "", fmt.Errorf("no resources found in namespace %q", namespace)
	}

	// Build archive path in current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("backup-%s-%s.tar.gz", namespace, timestamp)
	outputPath := filepath.Join(wd, filename)

	files := make([]File, 0, len(manifests))
	for _, m := range manifests {
		relPath := filepath.Join(namespace, m.Filename)
		files = append(files, File{
			Name: relPath,
			Data: m.Content,
		})
	}

	if err := CreateArchive(outputPath, files); err != nil {
		return "", fmt.Errorf("create archive: %w", err)
	}

	return outputPath, nil
}

// RestoreNamespace restores resources from a tar.gz archive into the cluster.
// If namespaceOverride is non-empty, it is used as a default namespace for
// namespaceless manifests.
func RestoreNamespace(archivePath, kubeconfigPath, namespaceOverride string, cfg *rest.Config) error {
	if archivePath == "" {
		return fmt.Errorf("archive path is required")
	}

	var (
		client *k8s.Client
		err    error
	)

	if kubeconfigPath != "" {
		client, err = k8s.NewClient(kubeconfigPath)
	} else {
		client, err = k8s.NewClientFromDefault()
	}
	if err != nil {
		return fmt.Errorf("create Kubernetes client: %w", err)
	}

	files, err := ExtractArchive(archivePath)
	if err != nil {
		return fmt.Errorf("extract archive: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("archive %q is empty", archivePath)
	}

	// Prepare dynamic client and RESTMapper based on provided REST config.
	if cfg == nil {
		return fmt.Errorf("REST config is required")
	}
	disco, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return fmt.Errorf("create discovery client: %w", err)
	}
	cachedDisco := memory.NewMemCacheClient(disco)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDisco)
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("create dynamic client: %w", err)
	}

	ctx := context.Background()
	for _, f := range files {
		if err := client.ApplyYAML(ctx, mapper, dyn, namespaceOverride, f.Data); err != nil {
			return fmt.Errorf("apply manifest %s: %w", f.Name, err)
		}
	}

	return nil
}
