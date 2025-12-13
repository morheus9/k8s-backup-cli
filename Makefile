# Makefile for k8s-backup-cli (POSIX compatible)
BINARY_NAME ?= kubectl-backup
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Directories
CMD_DIR = cmd
BIN_DIR = bin
DIST_DIR = dist
COVERAGE_DIR = coverage

# Go tools
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOVET = $(GOCMD) vet
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

# Пути к инструментам
GOBIN = $(shell go env GOPATH 2>/dev/null || echo $(HOME)/go)/bin
export PATH := $(GOBIN):$(PATH)

# LDFlags
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := build

##@ Development
.PHONY: build install deps tidy run ensure-binary

build: ## Build binary for current platform
	@echo "Building $(BINARY_NAME) $(VERSION) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BIN_DIR)
	@if $(GOBUILD) -trimpath $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR); then \
		chmod +x $(BIN_DIR)/$(BINARY_NAME); \
		echo "✅ Binary created: $(BIN_DIR)/$(BINARY_NAME)"; \
	else \
		echo "❌ Build failed"; \
		exit 1; \
	fi

ensure-binary: ## Ensure binary exists and is executable
	@if [ ! -f $(BIN_DIR)/$(BINARY_NAME) ] || [ ! -x $(BIN_DIR)/$(BINARY_NAME) ]; then \
		$(MAKE) build; \
	fi

install: build ## Install binary to system path
	@echo "Installing to /usr/local/bin..."
	@if [ -w /usr/local/bin ]; then \
		install -m 0755 $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/; \
	else \
		sudo install -m 0755 $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/; \
	fi
	@echo "✅ Installed! Run with: $(BINARY_NAME) --help"

deps: ## Download dependencies
	$(GOMOD) download
	@echo "✅ Dependencies downloaded"

tidy: ## Tidy go.mod
	$(GOMOD) tidy
	@echo "✅ Go modules tidied"

run: ensure-binary ## Run with arguments (make run ARGS="backup default --password test")
	@./$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

##@ Testing
.PHONY: test test-unit test-integration coverage

test: test-unit ## Run all tests
	@echo "✅ All tests passed!"

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./...

test-integration: ensure-binary ## Run integration tests (requires k8s cluster)
	@echo "Running integration tests..."
	@echo "Creating test namespace..."
	-kubectl create namespace backup-test 2>/dev/null || true
	@echo "Creating test resources..."
	-kubectl -n backup-test create configmap test-cm --from-literal=key=value 2>/dev/null || true
	-kubectl -n backup-test create secret generic test-secret --from-literal=password=123 2>/dev/null || true
	@echo "Running backup..."
	./$(BIN_DIR)/$(BINARY_NAME) backup backup-test --password test123 --output ./test-backups 2>&1
	@if ls ./test-backups/backup-*.tar.gz >/dev/null 2>&1; then \
		echo "✅ Backup created successfully!"; \
	else \
		echo "❌ Backup failed!"; exit 1; \
	fi
	@echo "Cleaning up..."
	rm -rf ./test-backups 2>/dev/null || true
	-kubectl delete namespace backup-test 2>/dev/null || true
	@echo "✅ Integration tests passed!"

coverage: ## Generate test coverage report
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "📊 Coverage report: $(COVERAGE_DIR)/coverage.html"

##@ Code Quality
.PHONY: fmt vet lint lint-fix security check

fmt: ## Format Go code
	$(GOFMT) ./...
	@echo "✅ Code formatted"

vet: ## Run go vet
	$(GOVET) ./...
	@echo "✅ Go vet passed"

lint: ## Run golint
	@if ! command -v golint >/dev/null 2>&1; then \
		echo "golint not found. Installing latest version..."; \
		go install golang.org/x/lint/golint@latest; \
	fi
	golint ./...
	@echo "✅ Lint passed"

lint-fix: ## Fix linting issues
	@echo "Note: golint doesn't support auto-fixing. Running gofmt and goimports instead..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	goimports -w .
	$(GOFMT) ./...
	@echo "✅ Formatting done"

security: ## Run security checks
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "gosec not found. Installing latest version..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...
	@echo "✅ Security checks passed"

check: fmt vet lint security ## Run all code quality checks
	@echo "✅ All checks passed!"

##@ Tool installation
.PHONY: tools

tools: ## Install all development tools
	@echo "Installing development tools..."
	go install golang.org/x/lint/golint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "✅ Tools installed"

##@ Build & Release
.PHONY: cross-build release clean distclean snapshot

cross-build: ## Build for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then \
				ext=".exe"; \
			else \
				ext=""; \
			fi; \
			echo "Building $$os/$$arch..."; \
			if GOOS=$$os GOARCH=$$arch $(GOBUILD) -trimpath $(LDFLAGS) \
				-o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext ./$(CMD_DIR); then \
				[ "$$os" != "windows" ] && chmod +x $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext || true; \
				echo "  ✅ Success"; \
			else \
				echo "  ❌ Failed"; \
			fi; \
		done; \
	done
	@echo "📦 Builds available in $(DIST_DIR)/"

release: cross-build checksums ## Create release package
	@echo "Creating release $(VERSION)..."
	cp README.md LICENSE $(DIST_DIR)/ 2>/dev/null || true
	@echo "🚀 Release $(VERSION) created in $(DIST_DIR)/"

checksums: ## Generate checksums for releases
	@echo "Generating checksums..."
	cd $(DIST_DIR) && \
		shasum -a 256 * > sha256sums.txt 2>/dev/null || \
		sha256sum * > sha256sums.txt 2>/dev/null || true
	@echo "✅ Checksums generated"

snapshot: ## Create development snapshot
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) -trimpath $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-snapshot ./$(CMD_DIR)
	chmod +x $(DIST_DIR)/$(BINARY_NAME)-snapshot
	@echo "📸 Snapshot: $(DIST_DIR)/$(BINARY_NAME)-snapshot"

##@ Utilities
.PHONY: help clean

clean: ## Clean build artifacts
	rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR) test-backups
	@echo "🧹 Clean completed"

distclean: clean ## Deep clean (includes dependencies)
	$(GOCMD) clean -cache -testcache -modcache
	@echo "🧹 Deep clean completed"

##@ Documentation
.PHONY: docs

docs: ensure-binary ## Generate CLI documentation
	@mkdir -p docs
	chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@if ./$(BIN_DIR)/$(BINARY_NAME) --help > docs/cli-reference.md 2>/dev/null; then \
		echo "✅ CLI reference generated"; \
	else \
		echo "⚠️  Could not generate CLI reference"; \
		echo "# CLI Reference\n\nCommand not available yet." > docs/cli-reference.md; \
	fi
	@if ./$(BIN_DIR)/$(BINARY_NAME) backup --help > docs/backup-command.md 2>/dev/null; then \
		echo "✅ Backup command help generated"; \
	else \
		echo "# Backup Command\n\nBackup command not available yet." > docs/backup-command.md; \
	fi
	@if ./$(BIN_DIR)/$(BINARY_NAME) restore --help > docs/restore-command.md 2>/dev/null; then \
		echo "✅ Restore command help generated"; \
	else \
		echo "# Restore Command\n\nRestore command not available yet." > docs/restore-command.md; \
	fi
	@echo "📚 Documentation generated in docs/"

##@ Help
.PHONY: help

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "Quick start:"
	@echo "  make tools      - Install development tools"
	@echo "  make build      - Build the binary"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make run ARGS='--help' - Run with arguments"
	@echo "  make docs       - Generate documentation"
	@echo "  make test       - Run tests"
	@echo "  make check      - Run all code quality checks"