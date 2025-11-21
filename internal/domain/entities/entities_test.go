package entities

import "testing"

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.2.0", "1.1.0", 1},
		{"1.1.5", "1.1.4", 1},
	}

	for _, tt := range tests {
		v1, _ := NewVersion(tt.v1)
		v2, _ := NewVersion(tt.v2)
		got := v1.Compare(v2)

		if got != tt.want {
			t.Errorf("Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, got, tt.want)
		}
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1.0.0", "1.0.0"},
		{"2.5.1", "2.5.1"},
		{"10.0.0", "10.0.0"},
	}

	for _, tt := range tests {
		v, err := NewVersion(tt.input)
		if err != nil {
			t.Fatalf("NewVersion(%q) error: %v", tt.input, err)
		}

		if got := v.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestRegistry_Operations(t *testing.T) {
	r := NewRegistry()

	if r.HasPackage("test") {
		t.Error("empty registry should not have packages")
	}

	v, _ := NewVersion("1.0.0")
	pkg := NewPackage("test", PackageTypeCLI, v)
	pkg.InstallPath = "/test/path"

	r.AddPackage(pkg)

	if !r.HasPackage("test") {
		t.Error("package should exist after adding")
	}

	if _, ok := r.GetPackage("test", "1.0.0"); !ok {
		t.Error("should find added package")
	}

	r.SetGlobalVersion("test", "1.0.0")
	if ver, ok := r.GetGlobalVersion("test"); !ok || ver != "1.0.0" {
		t.Errorf("global version = %q, want %q", ver, "1.0.0")
	}

	r.RemovePackage("test", "1.0.0")
	if r.HasPackage("test") {
		t.Error("package should be removed")
	}
}

func TestPackage_Identifier(t *testing.T) {
	v, _ := NewVersion("1.2.3")
	pkg := NewPackage("nano", PackageTypeCLI, v)

	want := "nano@1.2.3"
	if got := pkg.Identifier(); got != want {
		t.Errorf("Identifier() = %q, want %q", got, want)
	}
}

func TestFormula_GetPlatformConfig(t *testing.T) {
	f := NewFormula("test", PackageTypeCLI)

	if f.GetPlatformConfig("darwin", "arm64") != nil {
		t.Error("expected nil for empty platforms")
	}

	// Add platform config
	f.Platforms["darwin"] = ArchConfig{
		"arm64": &PlatformConfig{DownloadURL: "https://example.com"},
	}

	cfg := f.GetPlatformConfig("darwin", "arm64")
	if cfg == nil {
		t.Fatal("expected config for darwin/arm64")
	}

	if cfg.DownloadURL != "https://example.com" {
		t.Errorf("DownloadURL = %q, want %q", cfg.DownloadURL, "https://example.com")
	}
}

func TestFormula_TypeChecks(t *testing.T) {
	cli := NewFormula("cli-tool", PackageTypeCLI)
	if !cli.IsCLI() {
		t.Error("expected IsCLI() = true")
	}
	if cli.IsGUI() {
		t.Error("expected IsGUI() = false")
	}

	gui := NewFormula("gui-app", PackageTypeGUI)
	if !gui.IsGUI() {
		t.Error("expected IsGUI() = true")
	}
	if gui.IsCLI() {
		t.Error("expected IsCLI() = false")
	}
}

func TestWandRC_Operations(t *testing.T) {
	rc := NewWandRC()

	if rc.HasVersion("test") {
		t.Error("empty WandRC should not have versions")
	}

	rc.SetVersion("nano", "8.7.0")
	if !rc.HasVersion("nano") {
		t.Error("should have version after setting")
	}

	ver, ok := rc.GetVersion("nano")
	if !ok || ver != "8.7.0" {
		t.Errorf("GetVersion() = (%q, %v), want (%q, true)", ver, ok, "8.7.0")
	}

	rc.RemoveVersion("nano")
	if rc.HasVersion("nano") {
		t.Error("version should be removed")
	}
}
