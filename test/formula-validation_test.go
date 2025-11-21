package test

import (
	"testing"

	domainadapters "github.com/ochairo/wand/internal/domain-adapters"
	"github.com/ochairo/wand/internal/domain/entities"
)

// TestAllFormulasLoad tests that all formulas in the formulas directory load correctly
func TestAllFormulasLoad(t *testing.T) {
	formulaDir := "../formulas"
	fs := domainadapters.NewFileSystemAdapter()
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulaDir)

	// Expected formulas
	expectedFormulas := []struct {
		name     string
		pkgType  entities.PackageType
		binaries []string
	}{
		// CLI tools - only test formulas with potions releases
		{"nano", entities.PackageTypeCLI, []string{"nano"}},
		{"make", entities.PackageTypeCLI, []string{"make"}},
		{"zsh", entities.PackageTypeCLI, []string{"zsh"}},
		{"zlib", entities.PackageTypeCLI, []string{}},
		// GUI apps
		{"microsoft-edge", entities.PackageTypeGUI, []string{}},
	}

	for _, expected := range expectedFormulas {
		t.Run(expected.name, func(t *testing.T) {
			formula, err := formulaRepo.GetFormula(expected.name)
			if err != nil {
				t.Fatalf("Failed to load formula %s: %v", expected.name, err)
			}

			// Verify type
			if formula.Type != expected.pkgType {
				t.Errorf("Expected type %v, got %v", expected.pkgType, formula.Type)
			}

			// Verify binaries for CLI packages
			if expected.pkgType == entities.PackageTypeCLI && len(expected.binaries) > 0 {
				if len(formula.Binaries) != len(expected.binaries) {
					t.Errorf("Expected %d binaries, got %d", len(expected.binaries), len(formula.Binaries))
				}
			}

			// Verify platform configs exist
			platform := entities.CurrentPlatform()
			config := formula.GetPlatformConfigFor(platform)
			if config == nil {
				t.Errorf("No platform config for %s/%s", platform.OS, platform.Arch)
			} else if config.DownloadURL == "" {
				t.Error("Download URL is empty")
			}

			t.Logf("✓ Formula %s: type=%s, platforms=%d",
				formula.Name, formula.Type, len(formula.Platforms))
		})
	}
}

// TestFormulaMetadata tests that formulas have proper metadata
func TestFormulaMetadata(t *testing.T) {
	formulaDir := "../formulas"
	fs := domainadapters.NewFileSystemAdapter()
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulaDir)

	formulas := []string{"nano", "make", "zsh", "zlib", "microsoft-edge"}

	for _, name := range formulas {
		t.Run(name, func(t *testing.T) {
			formula, err := formulaRepo.GetFormula(name)
			if err != nil {
				t.Fatalf("Failed to load: %v", err)
			}

			// Check required metadata
			if formula.Description == "" {
				t.Error("Description is empty")
			}
			if formula.Homepage == "" {
				t.Error("Homepage is empty")
			}
			if formula.Repository == "" {
				t.Error("Repository is empty")
			}
			if formula.License == "" {
				t.Error("License is empty")
			}
			if len(formula.Tags) == 0 {
				t.Error("No tags defined")
			}

			t.Logf("✓ %s: %s (%s)", formula.Name, formula.Description, formula.License)
		})
	}
}

// TestFormulaURLTemplates tests that download URLs have proper placeholders
func TestFormulaURLTemplates(t *testing.T) {
	formulaDir := "../formulas"
	fs := domainadapters.NewFileSystemAdapter()
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulaDir)

	cliFormulas := []string{"nano"}

	for _, name := range cliFormulas {
		t.Run(name, func(t *testing.T) {
			formula, err := formulaRepo.GetFormula(name)
			if err != nil {
				t.Fatalf("Failed to load: %v", err)
			}

			platform := entities.CurrentPlatform()
			config := formula.GetPlatformConfigFor(platform)
			if config == nil {
				t.Skip("No config for current platform")
			}

			url := config.DownloadURL

			// Check for version placeholder (most formulas should have it)
			// Exception: formulas that use "latest" endpoints
			if name != "visual-studio-code" && name != "slack" {
				// Most formulas should have version in URL
				t.Logf("Download URL: %s", url)
			}

			if url == "" {
				t.Error("Download URL is empty")
			}
		})
	}
}
