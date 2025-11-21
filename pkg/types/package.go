// Package types provides public types for wand API consumers.
package types

import "time"

// PackageType represents the type of package
type PackageType string

const (
	// PackageTypeCLI represents a command-line tool package.
	PackageTypeCLI PackageType = "cli"
	// PackageTypeGUI represents a GUI application package.
	PackageTypeGUI PackageType = "gui"
	// PackageTypeDotfile represents a dotfile package.
	PackageTypeDotfile PackageType = "dotfile"
)

// Package represents an installed package
type Package struct {
	Name        string      // Package name (e.g., "node", "go", "nvim")
	Type        PackageType // CLI, GUI, or Dotfile
	Version     *Version    // Installed version
	InstalledAt time.Time   // Installation timestamp
	BinPath     string      // Path to binary/executable
	InstallPath string      // Full installation directory path
	IsGlobal    bool        // Whether this is the global version
}

// Identifier returns a unique identifier for the package
func (p *Package) Identifier() string {
	return p.Name + "@" + p.Version.String()
}

// VersionString returns the version as a string
func (p *Package) VersionString() string {
	if p.Version == nil {
		return ""
	}
	return p.Version.String()
}
