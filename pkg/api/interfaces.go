// Package api provides public interfaces for wand operations.
package api

import "github.com/ochairo/wand/pkg/types"

// Installer provides package installation operations.
type Installer interface {
	// Install installs a package with the specified version.
	// If version is empty or "latest", installs the latest version.
	Install(packageName, version string) (*types.Package, error)

	// Uninstall removes a specific version of a package.
	Uninstall(packageName, version string) error

	// IsInstalled checks if a specific package version is installed.
	IsInstalled(packageName, version string) bool

	// ListInstalled returns all installed versions of a package.
	ListInstalled(packageName string) ([]*types.Package, error)
}

// RegistryManager provides registry management operations.
type RegistryManager interface {
	// GetRegistry returns the current registry state.
	GetRegistry() (*types.Registry, error)

	// SetGlobalVersion sets the active global version for a package.
	SetGlobalVersion(packageName, version string) error

	// GetGlobalVersion returns the active global version for a package.
	GetGlobalVersion(packageName string) (string, error)

	// ListPackages returns all registered packages.
	ListPackages() ([]*types.PackageEntry, error)
}

// FormulaProvider provides formula lookup operations.
type FormulaProvider interface {
	// GetFormula retrieves a formula by name.
	GetFormula(name string) (*types.Formula, error)

	// ListFormulas returns all available formulas.
	ListFormulas() ([]*types.Formula, error)

	// SearchFormulas searches formulas by name or tag.
	SearchFormulas(query string) ([]*types.Formula, error)
}

// VersionResolver provides version resolution operations.
type VersionResolver interface {
	// ListAvailableVersions returns all available versions for a package.
	ListAvailableVersions(packageName string) ([]*types.Version, error)

	// ResolveVersion resolves a version constraint to a specific version.
	// Examples: "latest", "1.2.3", "^1.2.0", "~1.2.0"
	ResolveVersion(packageName, constraint string) (*types.Version, error)
}
