# autobackup_manifests_CLI

Stack:
- controller-gen
GOBIN=/home/pi/go/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
- kubebuilder
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
- kustomize
sudo apt install kustomize



kubebuilder init --domain example.com --repo github.com/morheus9/autobackup-manifest-operator
kubebuilder create api --group backup --version v1alpha1 --kind BackupSchedule