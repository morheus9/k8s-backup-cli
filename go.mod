module github.com/morheus9/k8s-backup-cli

go 1.25

require (
	// Для работы с Kubernetes
	k8s.io/apimachinery v0.29.0
	k8s.io/client-go v0.29.0

	// Для CLI интерфейса
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2

	// Для шифрования
	golang.org/x/crypto v0.17.0

	// Для работы с YAML/JSON
	gopkg.in/yaml.v3 v3.0.1

	// Для красивого вывода
	github.com/fatih/color v1.16.0
)
