package test

import (
	"os"
	"path/filepath"
	"testing"

	domainadapters "github.com/ochairo/wand/internal/domain-adapters"
	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/services"
)

// TestGUIWorkflow tests GUI application installation
func TestGUIWorkflow(t *testing.T) {
	// Create isolated test environment
	testHome := filepath.Join(os.TempDir(), "test-home-gui")
	wandDir := filepath.Join(testHome, ".wand")
	formulaDir := "../formulas"

	// Clean up before and after
	_ = os.RemoveAll(testHome)
	defer func() { _ = os.RemoveAll(testHome) }()

	// Create directories
	_ = os.MkdirAll(wandDir, 0755)                                 //nolint:gosec
	_ = os.MkdirAll(filepath.Join(testHome, "Applications"), 0755) //nolint:gosec

	// Initialize adapters
	fs := domainadapters.NewFileSystemAdapter()
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulaDir)
	_ = domainadapters.NewRegistryRepository(fs, wandDir)

	// Test: Load GUI formula
	t.Run("LoadGUIFormula", func(t *testing.T) {
		formula, err := formulaRepo.GetFormula("microsoft-edge")
		if err != nil {
			t.Fatalf("Failed to load GUI formula: %v", err)
		}

		if formula.Type != entities.PackageTypeGUI {
			t.Errorf("Expected type GUI, got %v", formula.Type)
		}

		if formula.AppName == "" {
			t.Error("AppName should not be empty for GUI formula")
		}

		t.Logf("✓ Loaded GUI formula: %s (type: %s, app: %s)", formula.Name, formula.Type, formula.AppName)
	})

	// Test: Verify platform config
	t.Run("VerifyPlatformConfig", func(t *testing.T) {
		formula, err := formulaRepo.GetFormula("microsoft-edge")
		if err != nil {
			t.Fatalf("Failed to load formula: %v", err)
		}

		platform := entities.CurrentPlatform()
		config := formula.GetPlatformConfigFor(platform)

		if config == nil {
			t.Fatalf("No platform config for %s/%s", platform.OS, platform.Arch)
		}

		if config.DownloadURL == "" {
			t.Error("Download URL should not be empty")
		}

		t.Logf("✓ Platform config found for %s/%s", platform.OS, platform.Arch)
		t.Logf("  Download URL: %s", config.DownloadURL)
	})

	// Test: Check registry integration
	t.Run("RegistryIntegration", func(t *testing.T) {
		registry := entities.NewRegistry()

		// Simulate GUI app installation (GUI apps use "0.0.0" as version placeholder)
		version, err := entities.NewVersion("0.0.0")
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}

		pkg := entities.NewPackage("visual-studio-code", entities.PackageTypeGUI, version)
		pkg.InstallPath = filepath.Join(wandDir, "apps", "visual-studio-code")
		pkg.IsGlobal = true

		registry.AddPackage(pkg)
		registry.SetGlobalVersion("visual-studio-code", "0.0.0")

		// Verify it's in registry
		if !registry.HasPackage("visual-studio-code") {
			t.Error("Package should be in registry")
		}

		installedPkg, exists := registry.GetPackage("visual-studio-code", "0.0.0")
		if !exists {
			t.Fatal("Should find package in registry")
		}

		if installedPkg.Type != entities.PackageTypeGUI {
			t.Errorf("Expected type GUI, got %v", installedPkg.Type)
		}

		t.Logf("✓ Registry integration works for GUI apps")
	})

	// Test: Multiple GUI apps
	t.Run("MultipleGUIApps", func(t *testing.T) {
		guiApps := []string{"microsoft-edge"}

		for _, appName := range guiApps {
			formula, err := formulaRepo.GetFormula(appName)
			if err != nil {
				t.Errorf("Failed to load %s: %v", appName, err)
				continue
			}

			if formula.Type != entities.PackageTypeGUI {
				t.Errorf("%s should be GUI type, got %v", appName, formula.Type)
			}

			if formula.AppName == "" {
				t.Errorf("%s should have app_name defined", appName)
			}

			t.Logf("✓ GUI app formula: %s -> %s", appName, formula.AppName)
		}
	})
}

// TestDMGExtraction tests DMG file handling
func TestDMGExtraction(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping DMG test in CI environment")
	}

	platform := entities.CurrentPlatform()
	if !platform.IsDarwin() {
		t.Skip("DMG extraction only supported on macOS")
	}

	t.Run("DMGSupportAvailable", func(t *testing.T) {
		fs := domainadapters.NewFileSystemAdapter()
		extractor := domainadapters.NewExtractorAdapter(fs)

		// Just verify the extractor interface exists
		// Actual DMG extraction would require a real DMG file
		if extractor == nil {
			t.Fatal("Extractor should not be nil")
		}

		t.Logf("✓ DMG extractor available for macOS")
	})
}

// TestGUIServiceIntegration tests WandfileService with GUI apps
func TestGUIServiceIntegration(t *testing.T) {
	testHome := filepath.Join(os.TempDir(), "test-home-gui-service")
	wandDir := filepath.Join(testHome, ".wand")

	_ = os.RemoveAll(testHome)
	defer func() { _ = os.RemoveAll(testHome) }()
	_ = os.MkdirAll(wandDir, 0755) //nolint:gosec

	fs := domainadapters.NewFileSystemAdapter()
	registryRepo := domainadapters.NewRegistryRepository(fs, wandDir)

	t.Run("CheckGUIAppsInWandfile", func(t *testing.T) {
		// Create wandfile with GUI apps
		wandfile := entities.NewWandfile()
		wandfile.AddGUI("microsoft-edge")

		if len(wandfile.GUI) != 1 {
			t.Errorf("Expected 1 GUI app, got %d", len(wandfile.GUI))
		}

		t.Logf("✓ Wandfile supports GUI apps: %v", wandfile.GUI)
	})

	t.Run("WandfileServiceChecksGUI", func(t *testing.T) {
		// Initialize minimal service dependencies
		formulaDir := "../formulas"
		formulaRepo := domainadapters.NewFormulaRepository(fs, formulaDir)
		wandfileRepo := domainadapters.NewWandfileRepository(fs)
		dotfileRepo := domainadapters.NewDotfileRepository(fs, wandDir)
		downloader := domainadapters.NewDownloaderAdapter()
		extractor := domainadapters.NewExtractorAdapter(fs)
		shellExecutor := domainadapters.NewShellExecutorAdapter()

		// Note: Using nil for testing - actual GitHub integration tested elsewhere
		versionSvc := services.NewVersionService(nil, formulaRepo)
		installerSvc := services.NewInstallerService(
			formulaRepo,
			registryRepo,
			downloader,
			extractor,
			fs,
			shellExecutor,
			versionSvc,
			wandDir,
			testHome,
		)

		wandfileSvc := services.NewWandfileService(
			wandfileRepo,
			registryRepo,
			installerSvc,
			versionSvc,
			dotfileRepo,
			fs,
			shellExecutor,
			testHome,
		)

		// Create wandfile with GUI apps
		wandfile := entities.NewWandfile()
		wandfile.AddGUI("microsoft-edge")

		// Check should report missing (not installed)
		missing, err := wandfileSvc.Check(wandfile)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}

		if len(missing) != 1 {
			t.Errorf("Expected 1 missing app, got %d: %v", len(missing), missing)
		}

		t.Logf("✓ WandfileService correctly checks GUI apps")
		t.Logf("  Missing apps: %v", missing)
	})
}
