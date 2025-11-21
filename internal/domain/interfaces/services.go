// Package interfaces defines all domain interfaces following the dependency inversion principle.
// The package name "interfaces" is a standard Go pattern for abstract contracts in domain-driven design.
package interfaces

import "github.com/ochairo/wand/internal/domain/entities"

// PackageInstaller defines the interface for installing packages
type PackageInstaller interface {
	Install(formula *entities.Formula, version *entities.Version) (*entities.Package, error)
	Uninstall(packageName, version string) error
	IsInstalled(packageName, version string) bool
}

// ShimManager defines the interface for managing shims
type ShimManager interface {
	CreateShim(packageName string) error
	RemoveShim(packageName string) error
	UpdateShim(packageName string) error
	ListShims() ([]string, error)
	ShimExists(packageName string) bool
}

// DotfileManager defines the interface for managing dotfiles
type DotfileManager interface {
	Init(repoURL string) error
	Sync() error
	Map(target, source string) error
	Unmap(target string) error
	List() (map[string]string, error)
	Status() (string, error)
}

// WandfileManager defines the interface for managing wandfiles
type WandfileManager interface {
	Install(wandfile *entities.Wandfile) error
	Update() error
	Check(wandfile *entities.Wandfile) ([]string, error)
	Dump() (*entities.Wandfile, error)
}

// RegistryManager defines the interface for managing the registry
type RegistryManager interface {
	GetRegistry() (*entities.Registry, error)
	UpdateRegistry(registry *entities.Registry) error
	AddPackage(pkg *entities.Package) error
	RemovePackage(packageName, version string) error
	SetGlobalVersion(packageName, version string) error
	GetGlobalVersion(packageName string) (string, error)
}
