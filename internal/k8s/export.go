package k8s

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Manifest represents a single Kubernetes object serialized to YAML.
type Manifest struct {
	Filename string
	Content  []byte
}

// isSystemObject returns true for Kubernetes system-managed objects that
// should not be included in user backups (for example kube-root-ca.crt).
func isSystemObject(meta metav1.ObjectMeta) bool {
	name := meta.Name

	// Well-known auto-created ConfigMap present in every namespace.
	if name == "kube-root-ca.crt" || strings.HasPrefix(name, "kube-root-ca.") {
		return true
	}

	// Skip objects from purely system namespaces.
	switch meta.Namespace {
	case "kube-system", "kube-public", "kube-node-lease":
		return true
	}

	return false
}

// ExportNamespaceManifests returns YAML manifests for supported resources in the namespace.
// Each resource is encoded as a separate YAML document.
func (c *Client) ExportNamespaceManifests(ctx context.Context, namespace string) ([]Manifest, error) {
	var manifests []Manifest

	// ConfigMaps
	configMaps, err := c.Clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	for _, cm := range configMaps.Items {
		if isSystemObject(cm.ObjectMeta) {
			continue
		}
		cm.APIVersion = "v1"
		cm.Kind = "ConfigMap"
		data, err := yaml.Marshal(cm)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ConfigMap %s: %w", cm.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("configmap-%s.yaml", cm.Name),
			Content:  data,
		})
	}

	// Secrets
	secrets, err := c.Clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Secrets: %w", err)
	}
	for _, secret := range secrets.Items {
		if isSystemObject(secret.ObjectMeta) {
			continue
		}
		secret.APIVersion = "v1"
		secret.Kind = "Secret"
		data, err := yaml.Marshal(secret)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Secret %s: %w", secret.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("secret-%s.yaml", secret.Name),
			Content:  data,
		})
	}

	// Services
	services, err := c.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Services: %w", err)
	}
	for _, svc := range services.Items {
		if isSystemObject(svc.ObjectMeta) {
			continue
		}
		svc.APIVersion = "v1"
		svc.Kind = "Service"
		data, err := yaml.Marshal(svc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Service %s: %w", svc.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("service-%s.yaml", svc.Name),
			Content:  data,
		})
	}

	// Deployments
	deployments, err := c.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	for _, deploy := range deployments.Items {
		if isSystemObject(deploy.ObjectMeta) {
			continue
		}
		deploy.APIVersion = "apps/v1"
		deploy.Kind = "Deployment"
		data, err := yaml.Marshal(deploy)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Deployment %s: %w", deploy.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("deployment-%s.yaml", deploy.Name),
			Content:  data,
		})
	}

	// StatefulSets
	statefulSets, err := c.Clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list StatefulSets: %w", err)
	}
	for _, sts := range statefulSets.Items {
		if isSystemObject(sts.ObjectMeta) {
			continue
		}
		sts.APIVersion = "apps/v1"
		sts.Kind = "StatefulSet"
		data, err := yaml.Marshal(sts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal StatefulSet %s: %w", sts.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("statefulset-%s.yaml", sts.Name),
			Content:  data,
		})
	}

	// DaemonSets
	daemonSets, err := c.Clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DaemonSets: %w", err)
	}
	for _, ds := range daemonSets.Items {
		if isSystemObject(ds.ObjectMeta) {
			continue
		}
		ds.APIVersion = "apps/v1"
		ds.Kind = "DaemonSet"
		data, err := yaml.Marshal(ds)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal DaemonSet %s: %w", ds.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("daemonset-%s.yaml", ds.Name),
			Content:  data,
		})
	}

	// Jobs
	jobs, err := c.Clientset.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Jobs: %w", err)
	}
	for _, job := range jobs.Items {
		if isSystemObject(job.ObjectMeta) {
			continue
		}
		job.APIVersion = "batch/v1"
		job.Kind = "Job"
		data, err := yaml.Marshal(job)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Job %s: %w", job.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("job-%s.yaml", job.Name),
			Content:  data,
		})
	}

	// CronJobs
	cronJobs, err := c.Clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CronJobs: %w", err)
	}
	for _, cj := range cronJobs.Items {
		if isSystemObject(cj.ObjectMeta) {
			continue
		}
		cj.APIVersion = "batch/v1"
		cj.Kind = "CronJob"
		data, err := yaml.Marshal(cj)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal CronJob %s: %w", cj.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("cronjob-%s.yaml", cj.Name),
			Content:  data,
		})
	}

	// PersistentVolumeClaims
	pvcs, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PersistentVolumeClaims: %w", err)
	}
	for _, pvc := range pvcs.Items {
		if isSystemObject(pvc.ObjectMeta) {
			continue
		}
		pvc.APIVersion = "v1"
		pvc.Kind = "PersistentVolumeClaim"
		data, err := yaml.Marshal(pvc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal PersistentVolumeClaim %s: %w", pvc.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("pvc-%s.yaml", pvc.Name),
			Content:  data,
		})
	}

	// Ingresses
	ingresses, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Ingresses: %w", err)
	}
	for _, ing := range ingresses.Items {
		if isSystemObject(ing.ObjectMeta) {
			continue
		}
		ing.APIVersion = "networking.k8s.io/v1"
		ing.Kind = "Ingress"
		data, err := yaml.Marshal(ing)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Ingress %s: %w", ing.Name, err)
		}
		manifests = append(manifests, Manifest{
			Filename: fmt.Sprintf("ingress-%s.yaml", ing.Name),
			Content:  data,
		})
	}

	return manifests, nil
}
