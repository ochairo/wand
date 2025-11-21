# Development Guide

Technical documentation for Wand contributors.

## Table of Contents

- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Development Setup](#development-setup)
- [Building](#building)
- [Testing](#testing)
- [Cross-Platform Considerations](#cross-platform-considerations)
- [Code Standards](#code-standards)

## Architecture

Wand follows **Clean Architecture** principles with strict dependency inversion.

### Dependency Rule

```bash
┌─────────────────────────────────────────────────────────────┐
│                            cmd/                             │  ← Entry Point
│                     (CLI Interface)                         │
└──────────────────────────────┬──────────────────────────────┘
                               │
                           depends on
                               │
                               ↓
┌─────────────────────────────────────────────────────────────┐
│ ┌──────────────────────────┐    ┌─────────────────────────┐ │  ← Adapters
│ │     domain-adapters/     │    │   external-adapters/    │ │
│ │ (Domain interfaces impl) ├─//─┤ (External API clients)  │ │
│ └────────────┬─────────────┘    └────────────┬────────────┘ │
│              ↓                               ↓              │
│              └───────────────┬───────────────┘              │
└──────────────────────────────┼──────────────────────────────┘
                               │
                            depends on
                               │
                               ↓
┌─────────────────────────────────────────────────────────────┐
│                    domain-orchestrators/                    │  ← Use Cases
└──────────────────────────────┬──────────────────────────────┘
                               │
                            depends on
                               │
                               ↓
┌─────────────────────────────────────────────────────────────┐
│                           domain/                           │  ← Core
│                                                             │
│      ┌───────────────────────────────────────────────┐      │
│      │                  services/                    │      │
│      │               (Business Logic)                │      │
│      └───────┬──────────────────────────────┬────────┘      │
│              │                              │               │
│          depends on                        uses             │
│              │                              │               │
│              ↓                              ↓               │
│      ┌────────────────┐            ┌──────────────────┐     │
│      │  interfaces/   │ implements │    entities/     │     │
│      │  (Contracts)   │ ←──────────┤  (Core Objects)  │     │
│      └────────────────┘            └──────────────────┘     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### `cmd/wand/`

- Entry point and CLI interface
- Command parsing and user input validation
- Depends on: `domain-orchestrators/`
- Zero business logic - delegates to orchestrators

#### `internal/domain/`

Core business rules - no external dependencies

**`entities/`**

- Pure data structures
- Business objects with minimal behavior
- No dependencies outside domain

**`interfaces/`**

- Contracts that define what adapters must implement
- Repositories: Data persistence contracts
- Gateways: External service contracts
- Services: Business logic contracts

**`services/`**

- Core business logic
- Version comparison, validation
- Installation rules, switching rules
- Depends only on: `entities/` and `interfaces/`

#### `internal/domain-orchestrators/`

- Coordinates multiple domain services
- Implements use cases/workflows
- Depends on: `domain/services/` and `domain/interfaces/`
- Called by: `cmd/`

#### `internal/domain-adapters/`

- Implements domain interfaces
- Adapts external concerns to domain needs
- Depends on: `domain/interfaces/` and `external-adapters/`

**`repositories/`**

- Implements repository interfaces
- JSON, YAML, or database persistence

**`gateways/`**

- Implements gateway interfaces
- HTTP clients, filesystem operations, GitHub API

#### `internal/external-adapters/`

- Thin wrappers around third-party libraries
- Isolates external dependencies
- Makes libraries easier to mock/test

### Dependency Flow

```bash
cmd/wand
    ↓ uses
domain-orchestrators
    ↓ uses
domain/services ← implements ← domain/interfaces
    ↓ uses                           ↑
domain/entities                      │ implements
                                     │
                              domain-adapters
                                     ↓ uses
                              external-adapters
```

## Directory Structure

```bash
wand/
├── cmd/wand/                  # CLI entry point
├── internal/
│   ├── domain/                # Business logic (entities, services, interfaces)
│   ├── domain-orchestrators/  # Use cases (workflows)
│   ├── domain-adapters/       # Infrastructure (repositories, gateways)
│   └── external-adapters/     # External libraries (Cobra, etc)
├── docs/                      # Documentation
└── Makefile, go.mod, README.md
```

Runtime files created in `~/.wand/`:

```bash
~/.wand/
├── shims/          # Command shims
├── packages/       # Versioned packages (e.g., node/20.0.0/bin/)
├── apps/           # GUI applications
└── registry.json   # Package metadata
```

## Development Setup

### Prerequisites

- Go 1.23 or later
- Git
- Make (optional)

### Clone and Build

```bash
# Clone repository
git clone https://github.com/ochairo/wand.git
cd wand

# Install dependencies
go mod download

# Build
go build -o bin/wand ./cmd/wand

# Run
./bin/wand --help
```

### Using Makefile

```bash
# Build
make build

# Run tests
make test

# Install locally
make install

# Clean build artifacts
make clean
```

## Building

### Single Platform

```bash
# Build for current platform
go build -o bin/wand ./cmd/wand
```

### Cross-Platform

```bash
# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o bin/wand-darwin-amd64 ./cmd/wand

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o bin/wand-darwin-arm64 ./cmd/wand

# Linux x86_64
GOOS=linux GOARCH=amd64 go build -o bin/wand-linux-amd64 ./cmd/wand

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o bin/wand-linux-arm64 ./cmd/wand
```

### Build All Platforms

```bash
make build-all
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Specific Package

```bash
go test ./internal/domain/services/
```

### With Coverage

```bash
go test -cover ./...

# Detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Verbose Output

```bash
go test -v ./...
```

### Using Makefile

```bash
make test
make test-coverage
```

## Cross-Platform Considerations

Wand supports **macOS** and **Linux** on both **amd64** and **arm64** architectures.

### Supported Platforms

- **macOS (darwin):** amd64 (Intel), arm64 (Apple Silicon)
- **Linux:** amd64 (x86_64), arm64 (aarch64)

### Platform Detection

```go
import "runtime"

platform := entities.CurrentPlatform()
fmt.Printf("Running on: %s\n", platform.String()) // e.g., "darwin/arm64"

if platform.IsDarwin() {
    // macOS-specific code
}

if platform.IsLinux() {
    // Linux-specific code
}
```

### Platform-Specific Behavior

#### Shim Scripts

- **Both platforms:** Bash scripts in `~/.wand/shims/`
- **Portable:** Use `#!/usr/bin/env bash`
- **Compatibility:** POSIX-compliant scripts work everywhere

#### GUI Applications

- **macOS:**
  - `.app` bundles in `~/.wand/apps/`
  - Symlinked to `~/Applications/`

- **Linux:**
  - Binaries in `~/.wand/apps/{name}/`
  - Desktop entries in `~/.local/share/applications/`
  - Icons in appropriate locations

#### File Operations

- **Permissions:** Use POSIX permissions (755 for executables)
- **Symlinks:** Work identically on both platforms
- **Paths:** Use `filepath.Join()` for cross-platform path handling

### Testing Cross-Platform

Use Docker for Linux testing on macOS:

```bash
# Linux amd64
docker run --rm -it -v $(pwd):/work golang:1.23 bash
cd /work
go test ./...

# Linux arm64 (on Apple Silicon)
docker run --rm -it --platform linux/arm64 -v $(pwd):/work golang:1.23 bash
cd /work
go test ./...
```

## Code Standards

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golint` for linting
- Write godoc comments for exported functions

### Package Organization

- One package per directory
- Keep packages focused and cohesive
- Avoid circular dependencies

### Naming

- **Packages:** Short, lowercase, no underscores
- **Files:** Lowercase with underscores
- **Types:** PascalCase
- **Functions:** camelCase (unexported), PascalCase (exported)
- **Interfaces:** Often end with "-er" (e.g., `Installer`, `Repository`)

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to install package: %w", err)
}

// Check errors immediately
result, err := doSomething()
if err != nil {
    return err
}
```

### Testing

- Write tests for all business logic
- Use table-driven tests for multiple cases
- Mock external dependencies
- Aim for >80% code coverage

### Example Test

```go
func TestVersionCompare(t *testing.T) {
    tests := []struct {
        name     string
        v1       string
        v2       string
        expected int
    }{
        {"equal versions", "1.0.0", "1.0.0", 0},
        {"v1 greater", "2.0.0", "1.0.0", 1},
        {"v1 less", "1.0.0", "2.0.0", -1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CompareVersions(tt.v1, tt.v2)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

## Example Flows

### Installing a Package

1. User runs: `wand install node@20.0.0`

2. `cmd/wand/main.go`
   - Parses arguments (package: node, version: 20.0.0)
   - Calls `InstallCommandHandler`

3. `domain-orchestrators/command_handlers.go`
   - Validates package name and version
   - Calls `InstallOrchestrator`

4. `domain-orchestrators/install_orchestrator.go`
   - Loads formula via `FormulaRepository`
   - Checks version via `GitHubGateway`
   - Calls `InstallerService`

5. `domain/services/installer_service.go`
   - Downloads binary via `DownloadGateway`
   - Extracts to `~/.wand/packages/node/20.0.0/`
   - Creates shim via `ShimService`
   - Updates registry via `RegistryRepository`

6. Each layer only knows about the layer below it

### Shim Execution

1. User runs: `node --version` in `~/projects/my-app/`

2. Shell finds: `~/.wand/shims/node`

3. Shim script:
   - Walks up from `~/projects/my-app/`
   - Finds `~/projects/my-app/.wandrc`
   - Parses YAML: `node: 16.20.0`
   - Execs: `~/.wand/packages/node/16.20.0/bin/node --version`

4. Output: `v16.20.0`

## Dependencies

Core dependencies:

- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/google/go-github/v57` - GitHub API (planned)

## Release Process

1. Update version in `cmd/wand/version.go`
2. Update `CHANGELOG.md`
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. Build all platforms
6. Create GitHub release with binaries

## Performance Considerations

- Shim overhead should be <5ms
- Cache formula data locally
- Minimize filesystem operations
- Use goroutines for parallel downloads (when installing multiple packages)

## Security

- Verify checksums for downloads
- Use HTTPS only
- Validate symlink targets
- Sanitize user inputs
- No shell injection in shims

## Contributing Formulas

Formulas define how Wand installs software and are stored in the `ochairo/potions` repository.

### Quick Example

```yaml
name: nano
type: cli
description: Text editor
homepage: https://www.nano-editor.org
repository: ochairo/formulas
license: GPLv3

binaries:
  - nano

platforms:
  darwin:
    amd64:
      download_url: "https://releases.example.com/{version}/nano-darwin-amd64.tar.gz"
      checksum_url: "https://releases.example.com/{version}/nano-darwin-amd64.tar.gz.sha256"
    arm64:
      download_url: "https://releases.example.com/{version}/nano-darwin-arm64.tar.gz"
      checksum_url: "https://releases.example.com/{version}/nano-darwin-arm64.tar.gz.sha256"

  linux:
    amd64:
      download_url: "https://releases.example.com/{version}/nano-linux-amd64.tar.gz"
      checksum_url: "https://releases.example.com/{version}/nano-linux-amd64.tar.gz.sha256"
    arm64:
      download_url: "https://releases.example.com/{version}/nano-linux-arm64.tar.gz"
      checksum_url: "https://releases.example.com/{version}/nano-linux-arm64.tar.gz.sha256"
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique package identifier |
| `type` | string | `cli` or `gui` |
| `description` | string | Short description |
| `homepage` | string | Official website |
| `repository` | string | GitHub repo (owner/repo) |
| `platforms` | map | Platform-specific downloads |

### Schema

```yaml
name: string                       # Package name
type: string                       # cli or gui
description: string
homepage: string
repository: string
license: string                    # Optional
tags: [string]                     # Optional
binaries: [string]                 # For CLI packages
bin_path: string                   # Optional - defaults to root
app_name: string                   # For GUI packages (macOS)

platforms:
  darwin:
    amd64:
      download_url: string         # URL with {version} placeholder
      checksum_url: string         # Optional - verifies download
    arm64:
      download_url: string
      checksum_url: string
  linux:
    amd64:
      download_url: string
      checksum_url: string
    arm64:
      download_url: string
      checksum_url: string
```

### Validation

```bash
wand validate my-formula.yaml
```

### Submit Formula

1. Create `.yaml` file following the schema
2. Test: `wand install mypackage@1.0.0`
3. Submit pull request to `ochairo/potions`

### Local Formulas

Create custom formulas in `~/.wand/formulas/`:

```bash
mkdir -p ~/.wand/formulas
# Add your .yaml files here
wand search my-custom-package
```

For details, see the [potions repository](https://github.com/ochairo/potions).

## Contributing to Wand Core

See [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Git workflow
- Commit conventions
- Pull request process
- Code review guidelines
