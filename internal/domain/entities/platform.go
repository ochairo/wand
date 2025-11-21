package entities

import (
	"runtime"
	"strings"
)

// Platform represents the current operating system and architecture
type Platform struct {
	OS   string // darwin, linux
	Arch string // amd64, arm64
}

// CurrentPlatform returns the current platform
func CurrentPlatform() *Platform {
	return &Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// String returns the platform as "os/arch"
func (p *Platform) String() string {
	return p.OS + "/" + p.Arch
}

// IsDarwin returns true if running on macOS
func (p *Platform) IsDarwin() bool {
	return p.OS == "darwin"
}

// IsLinux returns true if running on Linux
func (p *Platform) IsLinux() bool {
	return p.OS == "linux"
}

// IsARM64 returns true if running on ARM64 architecture
func (p *Platform) IsARM64() bool {
	return p.Arch == "arm64"
}

// IsAMD64 returns true if running on AMD64 architecture
func (p *Platform) IsAMD64() bool {
	return p.Arch == "amd64"
}

// IsSupported returns true if the platform is supported
func (p *Platform) IsSupported() bool {
	supportedOS := p.IsDarwin() || p.IsLinux()
	supportedArch := p.IsAMD64() || p.IsARM64()
	return supportedOS && supportedArch
}

// NormalizeOS normalizes OS names for formula matching
// Handles variations like "macos" -> "darwin"
func NormalizeOS(os string) string {
	os = strings.ToLower(os)
	switch os {
	case "macos", "osx", "mac":
		return "darwin"
	case "linux":
		return "linux"
	default:
		return os
	}
}

// NormalizeArch normalizes architecture names
// Handles variations like "x86_64" -> "amd64"
func NormalizeArch(arch string) string {
	arch = strings.ToLower(arch)
	switch arch {
	case "x86_64", "x64", "amd64":
		return "amd64"
	case "aarch64", "arm64":
		return "arm64"
	default:
		return arch
	}
}
