// Package interfaces defines all domain interfaces following the dependency inversion principle.
// The package name "interfaces" is a standard Go pattern for abstract contracts in domain-driven design.
// See: https://github.com/golang/go/wiki/CodeReviewComments#package-names
package interfaces

import "github.com/ochairo/wand/internal/domain/entities"

// RegistryRepository defines the interface for registry persistence
type RegistryRepository interface {
	Load() (*entities.Registry, error)
	Save(registry *entities.Registry) error
	Exists() bool
}

// FormulaRepository defines the interface for formula retrieval
type FormulaRepository interface {
	GetFormula(name string) (*entities.Formula, error)
	ListFormulas() ([]*entities.Formula, error)
	Sync() error
}

// WandRCRepository defines the interface for .wandrc file operations
type WandRCRepository interface {
	Load(dir string) (*entities.WandRC, error)
	Save(dir string, wandrc *entities.WandRC) error
	Exists(dir string) bool
	FindInPath(startDir string) (*entities.WandRC, string, error)
}

// WandfileRepository defines the interface for wandfile operations
type WandfileRepository interface {
	Load(path string) (*entities.Wandfile, error)
	Save(path string, wandfile *entities.Wandfile) error
	Exists(path string) bool
}

// DotfileRepository defines the interface for dotfile operations
type DotfileRepository interface {
	Load() (*entities.DotfileConfig, error)
	Save(config *entities.DotfileConfig) error
	Exists() bool
}
