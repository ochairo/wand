package entities

import (
	"runtime"
	"testing"
)

func TestCurrentPlatform(t *testing.T) {
	platform := CurrentPlatform()

	if platform == nil {
		t.Fatal("CurrentPlatform() returned nil")
	}

	if platform.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, platform.OS)
	}

	if platform.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %s, got %s", runtime.GOARCH, platform.Arch)
	}
}

func TestPlatformString(t *testing.T) {
	tests := []struct {
		os   string
		arch string
		want string
	}{
		{"darwin", "amd64", "darwin/amd64"},
		{"darwin", "arm64", "darwin/arm64"},
		{"linux", "amd64", "linux/amd64"},
		{"linux", "arm64", "linux/arm64"},
	}

	for _, tt := range tests {
		p := &Platform{OS: tt.os, Arch: tt.arch}
		if got := p.String(); got != tt.want {
			t.Errorf("Platform{%s, %s}.String() = %s, want %s", tt.os, tt.arch, got, tt.want)
		}
	}
}

func TestPlatformChecks(t *testing.T) {
	darwin := &Platform{OS: "darwin", Arch: "arm64"}
	if !darwin.IsDarwin() {
		t.Error("Expected IsDarwin() to be true")
	}
	if darwin.IsLinux() {
		t.Error("Expected IsLinux() to be false")
	}
	if !darwin.IsARM64() {
		t.Error("Expected IsARM64() to be true")
	}

	linux := &Platform{OS: "linux", Arch: "amd64"}
	if !linux.IsLinux() {
		t.Error("Expected IsLinux() to be true")
	}
	if linux.IsDarwin() {
		t.Error("Expected IsDarwin() to be false")
	}
	if !linux.IsAMD64() {
		t.Error("Expected IsAMD64() to be true")
	}
}

func TestPlatformIsSupported(t *testing.T) {
	tests := []struct {
		name      string
		platform  *Platform
		supported bool
	}{
		{"macOS Intel", &Platform{"darwin", "amd64"}, true},
		{"macOS ARM", &Platform{"darwin", "arm64"}, true},
		{"Linux x64", &Platform{"linux", "amd64"}, true},
		{"Linux ARM", &Platform{"linux", "arm64"}, true},
		{"Windows x64", &Platform{"windows", "amd64"}, false},
		{"FreeBSD", &Platform{"freebsd", "amd64"}, false},
		{"32-bit", &Platform{"linux", "386"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.IsSupported(); got != tt.supported {
				t.Errorf("IsSupported() = %v, want %v", got, tt.supported)
			}
		})
	}
}

func TestNormalizeOS(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"darwin", "darwin"},
		{"macos", "darwin"},
		{"osx", "darwin"},
		{"mac", "darwin"},
		{"Darwin", "darwin"},
		{"MacOS", "darwin"},
		{"linux", "linux"},
		{"Linux", "linux"},
		{"windows", "windows"},
	}

	for _, tt := range tests {
		if got := NormalizeOS(tt.input); got != tt.want {
			t.Errorf("NormalizeOS(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeArch(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"amd64", "amd64"},
		{"x86_64", "amd64"},
		{"x64", "amd64"},
		{"AMD64", "amd64"},
		{"X86_64", "amd64"},
		{"arm64", "arm64"},
		{"aarch64", "arm64"},
		{"ARM64", "arm64"},
		{"386", "386"},
	}

	for _, tt := range tests {
		if got := NormalizeArch(tt.input); got != tt.want {
			t.Errorf("NormalizeArch(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
