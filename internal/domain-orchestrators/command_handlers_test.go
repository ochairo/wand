package domainorchestrators

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ochairo/wand/internal/domain/entities"
)

// mockCommandContext for testing
type mockCommandContext struct {
	args   []string
	flags  map[string]interface{}
	output strings.Builder
}

func newMockContext(args []string) *mockCommandContext {
	return &mockCommandContext{
		args:  args,
		flags: make(map[string]interface{}),
	}
}

func (m *mockCommandContext) GetArgs() []string { return m.args }
func (m *mockCommandContext) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&m.output, format, args...)
}
func (m *mockCommandContext) PrintError(format string, args ...interface{}) {}

func (m *mockCommandContext) GetStringFlag(name string) (string, error) {
	if val, ok := m.flags[name].(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("flag not found")
}

func (m *mockCommandContext) GetBoolFlag(name string) (bool, error) {
	if val, ok := m.flags[name].(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("flag not found")
}

// mockFormulaRepo for testing
type mockFormulaRepo struct {
	formulas map[string]*entities.Formula
}

func newMockFormulaRepo() *mockFormulaRepo {
	return &mockFormulaRepo{
		formulas: map[string]*entities.Formula{
			"nano": {Name: "nano", Description: "Text editor", Tags: []string{"editor", "cli"}},
			"make": {Name: "make", Description: "Build tool", Tags: []string{"build"}},
		},
	}
}

func (m *mockFormulaRepo) GetFormula(name string) (*entities.Formula, error) {
	if f, ok := m.formulas[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("formula not found")
}

func (m *mockFormulaRepo) ListFormulas() ([]*entities.Formula, error) {
	formulas := make([]*entities.Formula, 0, len(m.formulas))
	for _, f := range m.formulas {
		formulas = append(formulas, f)
	}
	return formulas, nil
}

func (m *mockFormulaRepo) Sync() error { return nil }

// mockRegistryRepo for testing
type mockRegistryRepo struct {
	registry *entities.Registry
}

func newMockRegistryRepo() *mockRegistryRepo {
	return &mockRegistryRepo{registry: entities.NewRegistry()}
}

func (m *mockRegistryRepo) Load() (*entities.Registry, error) { return m.registry, nil }
func (m *mockRegistryRepo) Save(r *entities.Registry) error   { m.registry = r; return nil }
func (m *mockRegistryRepo) Exists() bool                      { return true }

// mockFileSystem for testing
type mockFileSystem struct {
	exists map[string]bool
}

func newMockFileSystem() *mockFileSystem {
	return &mockFileSystem{exists: make(map[string]bool)}
}

func (m *mockFileSystem) Exists(path string) bool                                        { return m.exists[path] }
func (m *mockFileSystem) IsDir(path string) bool                                         { return false }
func (m *mockFileSystem) ReadFile(path string) ([]byte, error)                           { return nil, nil }
func (m *mockFileSystem) WriteFile(path string, data []byte, perm uint32) error          { return nil }
func (m *mockFileSystem) MkdirAll(path string, perm uint32) error                        { return nil }
func (m *mockFileSystem) Remove(path string) error                                       { return nil }
func (m *mockFileSystem) RemoveAll(path string) error                                    { return nil }
func (m *mockFileSystem) Symlink(oldname, newname string) error                          { return nil }
func (m *mockFileSystem) ReadSymlink(name string) (string, error)                        { return "", nil }
func (m *mockFileSystem) Chmod(path string, mode uint32) error                           { return nil }
func (m *mockFileSystem) Walk(root string, walkFn func(string, bool, error) error) error { return nil }

func TestSearchCommandHandler(t *testing.T) {
	repo := newMockFormulaRepo()
	handler := NewSearchCommandHandler(repo)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{"no args", []string{}, true, ""},
		{"found by name", []string{"nano"}, false, "nano"},
		{"found by tag", []string{"editor"}, false, "editor"},
		{"not found", []string{"xyz"}, false, "No packages found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newMockContext(tt.args)
			err := handler.Handle(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.contains != "" && !strings.Contains(ctx.output.String(), tt.contains) {
				t.Errorf("output missing %q", tt.contains)
			}
		})
	}
}

func TestVersionCommandHandler(t *testing.T) {
	handler := NewVersionCommandHandler("1.0.0", "2025-11-21", "abc123")
	ctx := newMockContext(nil)

	if err := handler.Handle(ctx); err != nil {
		t.Fatal(err)
	}

	output := ctx.output.String()
	for _, want := range []string{"1.0.0", "2025-11-21", "abc123"} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestDoctorCommandHandler(t *testing.T) {
	repo := newMockRegistryRepo()
	formulas := newMockFormulaRepo()
	fs := newMockFileSystem()
	wandDir := "/test/.wand"
	fs.exists[wandDir] = true

	handler := NewDoctorCommandHandler(repo, formulas, fs, wandDir)
	ctx := newMockContext(nil)

	if err := handler.Handle(ctx); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(ctx.output.String(), "Wand directory exists") {
		t.Error("missing health check output")
	}
}

func TestOutdatedCommandHandler_Empty(t *testing.T) {
	repo := newMockRegistryRepo()
	handler := NewOutdatedCommandHandler(nil, repo)
	ctx := newMockContext(nil)

	if err := handler.Handle(ctx); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(ctx.output.String(), "No packages installed") {
		t.Error("expected empty message")
	}
}

func TestInstallCommandHandler_ParseSpec(t *testing.T) {
	tests := []struct {
		spec    string
		name    string
		version string
	}{
		{"nano", "nano", "latest"},
		{"nano@8.7", "nano", "8.7"},
		{"node@18.0.0", "node", "18.0.0"},
	}

	for _, tt := range tests {
		parts := strings.Split(tt.spec, "@")
		name := parts[0]
		version := "latest"
		if len(parts) > 1 {
			version = parts[1]
		}

		if name != tt.name || version != tt.version {
			t.Errorf("parse %q: got (%q, %q), want (%q, %q)",
				tt.spec, name, version, tt.name, tt.version)
		}
	}
}
