package types

// Formula represents a package definition from the potions repository
type Formula struct {
	// Required fields
	Name        string
	Type        PackageType
	Description string
	Homepage    string
	Repository  string

	// Optional metadata
	License        string
	Tags           []string
	VersionPattern string

	// CLI-specific
	Binaries []string
	BinPath  string

	// GUI-specific (macOS)
	AppName string
	AppPath string

	// Version constraints
	MinVersion string
	MaxVersion string

	// Dependencies
	Dependencies []string
}
