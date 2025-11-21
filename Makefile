.PHONY: build test clean install run fmt lint help

# Binary name
BINARY_NAME=wand
INSTALL_PATH=/usr/local/bin

# Build variables
VERSION?=0.1.0-dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.Commit=${COMMIT}"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} ./cmd/wand

run: ## Run the application
	go run ./cmd/wand

install: build ## Install the binary to system
	@echo "Installing ${BINARY_NAME} to ${INSTALL_PATH}..."
	cp ${BINARY_NAME} ${INSTALL_PATH}/
	@echo "Installed successfully!"

uninstall: ## Uninstall the binary from system
	@echo "Removing ${BINARY_NAME} from ${INSTALL_PATH}..."
	rm -f ${INSTALL_PATH}/${BINARY_NAME}

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html
	rm -rf dist/

deps: ## Download dependencies
	go mod download
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

# Cross-platform build targets
PLATFORMS=darwin-amd64 darwin-arm64 linux-amd64 linux-arm64
dist: clean ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		OS=$${platform%-*}; \
		ARCH=$${platform#*-}; \
		OUTPUT="dist/${BINARY_NAME}-$$OS-$$ARCH"; \
		echo "Building $$platform..."; \
		GOOS=$$OS GOARCH=$$ARCH go build ${LDFLAGS} -o $$OUTPUT ./cmd/wand; \
	done
	@echo "Cross-platform builds complete in dist/"

dist-archive: dist ## Create release archives with checksums
	@echo "Creating release archives..."
	@cd dist && for platform in $(PLATFORMS); do \
		OS=$${platform%-*}; \
		ARCH=$${platform#*-}; \
		BINARY="${BINARY_NAME}-$$OS-$$ARCH"; \
		ARCHIVE="${BINARY_NAME}-${VERSION}-$$OS-$$ARCH.tar.gz"; \
		tar -czf $$ARCHIVE $$BINARY; \
		shasum -a 256 $$ARCHIVE > $$ARCHIVE.sha256; \
		echo "Created $$ARCHIVE"; \
	done
	@echo "Release archives created in dist/"

.DEFAULT_GOAL := help
