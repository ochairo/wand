package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	errs "github.com/ochairo/wand/internal/domain/errors"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// ShimService handles shim generation and version resolution
type ShimService struct {
	registryRepo interfaces.RegistryRepository
	wandrcRepo   interfaces.WandRCRepository
	formulaRepo  interfaces.FormulaRepository
	fs           interfaces.FileSystem
	wandDir      string
	shimTemplate string
}

// NewShimService creates a new shim service
func NewShimService(
	registryRepo interfaces.RegistryRepository,
	wandrcRepo interfaces.WandRCRepository,
	formulaRepo interfaces.FormulaRepository,
	fs interfaces.FileSystem,
	wandDir string,
) *ShimService {
	return &ShimService{
		registryRepo: registryRepo,
		wandrcRepo:   wandrcRepo,
		formulaRepo:  formulaRepo,
		fs:           fs,
		wandDir:      wandDir,
		shimTemplate: getShimTemplate(),
	}
}

// CreateShims creates shims for all binaries in a package
func (s *ShimService) CreateShims(packageName string, binaries []string) error {
	shimsDir := filepath.Join(s.wandDir, "shims")
	if err := s.fs.MkdirAll(shimsDir, 0755); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to create shims directory", err)
	}

	for _, binary := range binaries {
		shimPath := filepath.Join(shimsDir, binary)

		// Generate shim script
		shimContent := s.generateShimScript(binary, packageName)

		// Write shim file
		if err := s.fs.WriteFile(shimPath, []byte(shimContent), 0755); err != nil {
			return errs.Wrap(errs.ErrShimCreationFailed, fmt.Sprintf("Failed to create shim for %s", binary), err)
		}
	}

	return nil
}

// RemoveShims removes shims for a package's binaries
func (s *ShimService) RemoveShims(binaries []string) error {
	shimsDir := filepath.Join(s.wandDir, "shims")

	for _, binary := range binaries {
		shimPath := filepath.Join(shimsDir, binary)

		if s.fs.Exists(shimPath) {
			if err := s.fs.Remove(shimPath); err != nil {
				return errs.Wrap(errs.ErrShimExecutionFailed, fmt.Sprintf("Failed to remove shim for %s", binary), err)
			}
		}
	}

	return nil
}

// ResolveVersion resolves which version of a package to use
// Resolution order: .wandrc (project) -> global version from registry
func (s *ShimService) ResolveVersion(packageName, currentDir string) (string, error) {
	// Check for .wandrc in current directory or parent directories
	wandrc, wandrcPath, err := s.wandrcRepo.FindInPath(currentDir)
	if err == nil && wandrc != nil {
		// Check if package version is specified in .wandrc
		if version, exists := wandrc.Versions[packageName]; exists {
			return version, nil
		}

		// Continue searching in parent directory
		if wandrcPath != "" {
			parentDir := filepath.Dir(filepath.Dir(wandrcPath))
			if parentDir != "/" && parentDir != "." {
				return s.ResolveVersion(packageName, parentDir)
			}
		}
	}

	// Fall back to global version from registry
	registry, err := s.registryRepo.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load registry: %w", err)
	}

	globalVersion, exists := registry.GetGlobalVersion(packageName)
	if !exists {
		return "", fmt.Errorf("no version found for package %s", packageName)
	}

	return globalVersion, nil
}

// GetBinaryPath returns the full path to a package's binary
func (s *ShimService) GetBinaryPath(packageName, version, binaryName string) (string, error) {
	registry, err := s.registryRepo.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load registry: %w", err)
	}

	pkg, exists := registry.GetPackage(packageName, version)
	if !exists {
		return "", fmt.Errorf("package %s@%s not found in registry", packageName, version)
	}

	// Construct binary path
	binaryPath := filepath.Join(pkg.BinPath, binaryName)

	if !s.fs.Exists(binaryPath) {
		return "", fmt.Errorf("binary %s not found at %s", binaryName, binaryPath)
	}

	return binaryPath, nil
}

// generateShimScript generates the shim script content
func (s *ShimService) generateShimScript(binaryName, packageName string) string {
	script := s.shimTemplate
	script = strings.ReplaceAll(script, "{{WAND_DIR}}", s.wandDir)
	script = strings.ReplaceAll(script, "{{BINARY_NAME}}", binaryName)
	script = strings.ReplaceAll(script, "{{PACKAGE_NAME}}", packageName)
	return script
}

// getShimTemplate returns the shim script template
func getShimTemplate() string {
	return `#!/bin/sh
# Wand shim for {{BINARY_NAME}}
# This script automatically resolves the correct version based on .wandrc or global config

set -e

# Package info
PACKAGE_NAME="{{PACKAGE_NAME}}"
BINARY_NAME="{{BINARY_NAME}}"
WAND_DIR="{{WAND_DIR}}"

# Get current directory
CURRENT_DIR="$(pwd)"

# Function to find .wandrc and resolve version
resolve_version() {
    local dir="$1"
    while [ "$dir" != "/" ] && [ "$dir" != "." ]; do
        if [ -f "$dir/.wandrc" ]; then
            # Parse YAML to get version for this package
            # Looking for: versions:\n  package: version
            version=$(grep -A 100 "^versions:" "$dir/.wandrc" | grep "^  $PACKAGE_NAME:" | awk '{print $2}' | tr -d '"' | head -1)
            if [ -n "$version" ]; then
                echo "$version"
                return 0
            fi
        fi
        dir="$(dirname "$dir")"
    done

    # Fall back to global version from registry
    if [ -f "$WAND_DIR/registry.json" ]; then
        # Extract global version from JSON (simple jq alternative)
        version=$(grep -A 1 "\"$PACKAGE_NAME\":" "$WAND_DIR/registry.json" | tail -1 | tr -d ' ",')
        if [ -n "$version" ]; then
            echo "$version"
            return 0
        fi
    fi

    return 1
}

# Resolve version
VERSION=$(resolve_version "$CURRENT_DIR")

if [ -z "$VERSION" ]; then
    echo "Error: No version found for $PACKAGE_NAME" >&2
    echo "Run: wand install $PACKAGE_NAME" >&2
    exit 1
fi

# Construct binary path
BINARY_PATH="$WAND_DIR/packages/$PACKAGE_NAME/$VERSION/bin/$BINARY_NAME"

if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found at $BINARY_PATH" >&2
    echo "Try reinstalling: wand install $PACKAGE_NAME@$VERSION" >&2
    exit 1
fi

# Execute the actual binary with all arguments
exec "$BINARY_PATH" "$@"
`
}

// RefreshAllShims recreates all shims for all installed packages
func (s *ShimService) RefreshAllShims() error {
	registry, err := s.registryRepo.Load()
	if err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry", err)
	}

	// Remove all existing shims
	shimsDir := filepath.Join(s.wandDir, "shims")
	if s.fs.Exists(shimsDir) {
		if err := s.fs.RemoveAll(shimsDir); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to clear shims directory", err)
		}
	}

	// Recreate shims directory
	if err := s.fs.MkdirAll(shimsDir, 0755); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to create shims directory", err)
	}

	// Create shims for all CLI packages
	for _, entry := range registry.ListAllPackages() {
		if entry.Type != entities.PackageTypeCLI {
			continue
		}

		// Load formula to get actual binary names
		formula, err := s.formulaRepo.GetFormula(entry.Name)
		if err != nil {
			// Fallback: assume binary name matches package name
			formula = nil
		}

		// Get binaries from formula or use package name as fallback
		var binaries []string
		if formula != nil && len(formula.Binaries) > 0 {
			binaries = formula.Binaries
		} else {
			binaries = []string{entry.Name}
		}

		if err := s.CreateShims(entry.Name, binaries); err != nil {
			return errs.Wrap(errs.ErrShimCreationFailed, fmt.Sprintf("Failed to create shims for %s", entry.Name), err)
		}
	}

	return nil
}
