package k8s

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceInfo represents information about a Kubernetes resource
type ResourceInfo struct {
	Kind       string
	Name       string
	Namespace  string
	APIVersion string
}

func isSystemResource(meta metav1.ObjectMeta, kind string) bool {
	name := meta.Name

	// Well-known auto-created ConfigMap present in every namespace.
	if name == "kube-root-ca.crt" || strings.HasPrefix(name, "kube-root-ca.") {
		return true
	}

	// Cluster service in default namespace.
	if kind == "Service" && meta.Namespace == "default" && name == "kubernetes" {
		return true
	}

	// Purely system namespaces.
	switch meta.Namespace {
	case "kube-system", "kube-public", "kube-node-lease":
		return true
	}

	return false
}

// FetchResources fetches all resources in the specified namespace
func (c *Client) FetchResources(ctx context.Context, namespace string) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	// Fetch ConfigMaps
	configMaps, err := c.Clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	for _, cm := range configMaps.Items {
		if isSystemResource(cm.ObjectMeta, "ConfigMap") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "ConfigMap",
			Name:       cm.Name,
			Namespace:  cm.Namespace,
			APIVersion: "v1",
		})
	}

	// Fetch Secrets
	secrets, err := c.Clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Secrets: %w", err)
	}
	for _, secret := range secrets.Items {
		if isSystemResource(secret.ObjectMeta, "Secret") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "Secret",
			Name:       secret.Name,
			Namespace:  secret.Namespace,
			APIVersion: "v1",
		})
	}

	// Fetch Services
	services, err := c.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Services: %w", err)
	}
	for _, svc := range services.Items {
		if isSystemResource(svc.ObjectMeta, "Service") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "Service",
			Name:       svc.Name,
			Namespace:  svc.Namespace,
			APIVersion: "v1",
		})
	}

	// Fetch Deployments
	deployments, err := c.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	for _, deploy := range deployments.Items {
		if isSystemResource(deploy.ObjectMeta, "Deployment") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "Deployment",
			Name:       deploy.Name,
			Namespace:  deploy.Namespace,
			APIVersion: "apps/v1",
		})
	}

	// Fetch StatefulSets
	statefulSets, err := c.Clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list StatefulSets: %w", err)
	}
	for _, sts := range statefulSets.Items {
		if isSystemResource(sts.ObjectMeta, "StatefulSet") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "StatefulSet",
			Name:       sts.Name,
			Namespace:  sts.Namespace,
			APIVersion: "apps/v1",
		})
	}

	// Fetch DaemonSets
	daemonSets, err := c.Clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DaemonSets: %w", err)
	}
	for _, ds := range daemonSets.Items {
		if isSystemResource(ds.ObjectMeta, "DaemonSet") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "DaemonSet",
			Name:       ds.Name,
			Namespace:  ds.Namespace,
			APIVersion: "apps/v1",
		})
	}

	// Fetch Jobs
	jobs, err := c.Clientset.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Jobs: %w", err)
	}
	for _, job := range jobs.Items {
		if isSystemResource(job.ObjectMeta, "Job") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "Job",
			Name:       job.Name,
			Namespace:  job.Namespace,
			APIVersion: "batch/v1",
		})
	}

	// Fetch CronJobs
	cronJobs, err := c.Clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CronJobs: %w", err)
	}
	for _, cj := range cronJobs.Items {
		if isSystemResource(cj.ObjectMeta, "CronJob") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "CronJob",
			Name:       cj.Name,
			Namespace:  cj.Namespace,
			APIVersion: "batch/v1",
		})
	}

	// Fetch PersistentVolumeClaims
	pvcs, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PersistentVolumeClaims: %w", err)
	}
	for _, pvc := range pvcs.Items {
		if isSystemResource(pvc.ObjectMeta, "PersistentVolumeClaim") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "PersistentVolumeClaim",
			Name:       pvc.Name,
			Namespace:  pvc.Namespace,
			APIVersion: "v1",
		})
	}

	// Fetch Ingresses
	ingresses, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Ingresses: %w", err)
	}
	for _, ing := range ingresses.Items {
		if isSystemResource(ing.ObjectMeta, "Ingress") {
			continue
		}
		resources = append(resources, ResourceInfo{
			Kind:       "Ingress",
			Name:       ing.Name,
			Namespace:  ing.Namespace,
			APIVersion: "networking.k8s.io/v1",
		})
	}

	return resources, nil
}
