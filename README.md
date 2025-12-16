## Autobackup manifests CLI

### Install:

git clone https://github.com/morheus9/k8s-backup-cli.git

make install
```
Building kubectl-backup 87ff2a4-dirty for linux/amd64...
✅ Binary created: bin/kubectl-backup
Installing to /usr/local/bin...
✅ Installed! Run with: kubectl-backup --help
```

### Using:

kubectl-backup list your_namespace
```
KIND         NAME         NAMESPACE           API VERSION
----         ----         ---------           -----------
Secret       app-secret   your_namespace      v1
Deployment   my-app       your_namespace      apps/v1

Total: 2 resources
```

kubectl-backup backup your_namespace
```
Backup created at: /home/pi/Downloads/k8s-backup-cli/backup-your_namespace-20251215-210219.tar.gz
```

kubectl-backup restore your_namespace -f backup-your_namespace-20251215-210219.tar.gz
```
Successfully restored resources from backup-your_namespace-20251215-210219.tar.gz
```

### Uninstall:

make uninstall
```
Uninstalling kubectl-backup from /usr/local/bin...
✅ Uninstalled successfully
```
