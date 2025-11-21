package entities

// PlatformConfig represents platform-specific download configuration
type PlatformConfig struct {
	DownloadURL   string   `yaml:"download_url"`
	RequiresBuild bool     `yaml:"requires_build"`
	BuildCommands []string `yaml:"build_commands,omitempty"`
	ChecksumURL   string   `yaml:"checksum_url,omitempty"`
	DesktopFile   string   `yaml:"desktop_file,omitempty"` // Linux GUI
	IconFile      string   `yaml:"icon_file,omitempty"`    // Linux GUI
}

// ArchConfig maps architecture to platform configuration
type ArchConfig map[string]*PlatformConfig // amd64, arm64

// PlatformMap maps OS to architecture configurations
type PlatformMap map[string]ArchConfig // darwin, linux

// PostInstallHook defines post-installation commands
type PostInstallHook struct {
	Commands []string          `yaml:"commands,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
}

// Formula represents a package definition from the potions repository
type Formula struct {
	// Required fields
	Name        string      `yaml:"name"`
	Type        PackageType `yaml:"type"`
	Description string      `yaml:"description"`
	Homepage    string      `yaml:"homepage"`
	Repository  string      `yaml:"repository"`

	// Optional metadata
	License        string   `yaml:"license,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	VersionPattern string   `yaml:"version_pattern,omitempty"`

	// CLI-specific
	Binaries []string `yaml:"binaries,omitempty"`
	BinPath  string   `yaml:"bin_path,omitempty"`

	// GUI-specific (macOS)
	AppName string `yaml:"app_name,omitempty"`
	AppPath string `yaml:"app_path,omitempty"`

	// Platform downloads
	Platforms PlatformMap `yaml:"platforms"`

	// Hooks and dependencies
	PostInstall  *PostInstallHook `yaml:"post_install,omitempty"`
	Dependencies []string         `yaml:"dependencies,omitempty"`

	// Version constraints
	MinVersion string `yaml:"min_version,omitempty"`
	MaxVersion string `yaml:"max_version,omitempty"`
}

// NewFormula creates a new Formula
func NewFormula(name string, pkgType PackageType) *Formula {
	return &Formula{
		Name:      name,
		Type:      pkgType,
		Platforms: make(PlatformMap),
	}
}

// GetPlatformConfig returns platform configuration for OS/Arch
func (f *Formula) GetPlatformConfig(os, arch string) *PlatformConfig {
	if archConfig, ok := f.Platforms[os]; ok {
		if config, ok := archConfig[arch]; ok {
			return config
		}
	}
	return nil
}

// GetPlatformConfigFor returns platform configuration for a given platform
func (f *Formula) GetPlatformConfigFor(platform *Platform) *PlatformConfig {
	return f.GetPlatformConfig(platform.OS, platform.Arch)
}

// GetCurrentPlatformConfig returns the platform config for the current platform
func (f *Formula) GetCurrentPlatformConfig() *PlatformConfig {
	return f.GetPlatformConfigFor(CurrentPlatform())
}

// SupportsPlatform checks if a formula supports a given platform
func (f *Formula) SupportsPlatform(platform *Platform) bool {
	return f.GetPlatformConfigFor(platform) != nil
}

// SupportsCurrentPlatform checks if the formula supports the current platform
func (f *Formula) SupportsCurrentPlatform() bool {
	return f.SupportsPlatform(CurrentPlatform())
}

// GetDownloadURL returns the download URL for a platform and version
func (f *Formula) GetDownloadURL(platform *Platform, version string) string {
	config := f.GetPlatformConfigFor(platform)
	if config == nil || config.DownloadURL == "" {
		return ""
	}

	// Replace placeholders
	url := config.DownloadURL
	url = replacePlaceholder(url, "version", version)
	url = replacePlaceholder(url, "platform", platform.OS)
	url = replacePlaceholder(url, "arch", platform.Arch)

	return url
}

// GetCurrentDownloadURL returns the download URL for current platform and version
func (f *Formula) GetCurrentDownloadURL(version string) string {
	return f.GetDownloadURL(CurrentPlatform(), version)
}

// HasDependencies returns true if the formula has dependencies
func (f *Formula) HasDependencies() bool {
	return len(f.Dependencies) > 0
}

// HasPostInstall returns true if the formula has post-install hooks
func (f *Formula) HasPostInstall() bool {
	return f.PostInstall != nil && len(f.PostInstall.Commands) > 0
}

// replacePlaceholder replaces {key} with value in a string
func replacePlaceholder(s, key, value string) string {
	placeholder := "{" + key + "}"
	return replaceAll(s, placeholder, value)
}

// replaceAll is a simple string replace (avoiding strings import)
func replaceAll(s, old, replacement string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx == -1 {
			result += s
			break
		}
		result += s[:idx] + replacement
		s = s[idx+len(old):]
	}
	return result
}

// indexOf finds the index of substring in string
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// IsCLI returns true if this is a CLI package
func (f *Formula) IsCLI() bool {
	return f.Type == PackageTypeCLI
}

// IsGUI returns true if this is a GUI package
func (f *Formula) IsGUI() bool {
	return f.Type == PackageTypeGUI
}

// IsDotfile returns true if this is a dotfile package
func (f *Formula) IsDotfile() bool {
	return f.Type == PackageTypeDotfile
}
