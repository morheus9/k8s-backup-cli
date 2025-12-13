## Autobackup manifests operator




______________________________________________________________________
## For developer

### 1. Install utilities

- kubebuilder
```
go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

kubebuilder init --domain example.com --repo github.com/morheus9/autobackup-manifest-operator
kubebuilder create api --group backup --version v1alpha1 --kind BackupSchedule
```
- controller-gen
```
make controller-gen
GOBIN=/home/pi/go/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```
- kustomize
```
make kustomize
sudo apt install kustomize
```
### 2. Code generation

##### Генерирует код на основе маркеров в types.go
```
make generate
```
##### Generates CRD, RBAC, Webhook manifests
```
make manifests
```
### 3. Check build
```
make build
```
### 4. Installing in cluster
```
make install
```
### 5. Start in development mode
```
make run
```
##### Checking
```
kubectl get crd | grep backupschedule
```

make generate
make manifests
make run