package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/yaml"
)

// ApplyYAML applies a single Kubernetes manifest to the cluster.
// If the resource already exists, it will be updated.
func (c *Client) ApplyYAML(ctx context.Context, mapper *restmapper.DeferredDiscoveryRESTMapper, dyn dynamic.Interface, namespaceFallback string, manifest []byte) error {
	if len(manifest) == 0 {
		return nil
	}

	// Decode YAML into unstructured object.
	jsonData, err := yaml.YAMLToJSON(manifest)
	if err != nil {
		return fmt.Errorf("convert YAML to JSON: %w", err)
	}

	var obj unstructured.Unstructured
	if err := json.Unmarshal(jsonData, &obj); err != nil {
		return fmt.Errorf("unmarshal JSON into unstructured object: %w", err)
	}

	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("find REST mapping for %s: %w", gvk.String(), err)
	}

	ns := obj.GetNamespace()
	if ns == "" {
		ns = namespaceFallback
	}

	var resourceClient dynamic.ResourceInterface
	if mapping.Scope.Name() == "namespace" {
		resourceClient = dyn.Resource(mapping.Resource).Namespace(ns)
	} else {
		resourceClient = dyn.Resource(mapping.Resource)
	}

	// Try to create, fall back to update if already exists.
	// For create, resourceVersion must be empty.
	obj.SetResourceVersion("")
	_, err = resourceClient.Create(ctx, &obj, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		// Need current resource version for update.
		existing, getErr := resourceClient.Get(ctx, obj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("get existing %s/%s: %w", gvk.Kind, obj.GetName(), getErr)
		}
		obj.SetResourceVersion(existing.GetResourceVersion())
		if _, err := resourceClient.Update(ctx, &obj, metav1.UpdateOptions{}); err != nil {
			return fmt.Errorf("update %s/%s: %w", gvk.Kind, obj.GetName(), err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("create %s/%s: %w", gvk.Kind, obj.GetName(), err)
	}

	return nil
}
