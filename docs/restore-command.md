Restore Kubernetes resources from a previously created backup archive

Usage:
  kubectl-backup restore [flags]

Flags:
  -f, --file string         Path to backup archive (tar.gz) to restore from (required)
  -h, --help                help for restore
  -k, --kubeconfig string   Path to kubeconfig file (default: auto-detect)
  -n, --namespace string    Default namespace for namespaceless manifests
