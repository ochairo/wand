package test

import (
	"os"
	"path/filepath"
	"testing"

	domain_adapters "github.com/ochairo/wand/internal/domain-adapters"
	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/services"
	external_adapters "github.com/ochairo/wand/internal/external-adapters"
)

// TestWandfileWorkflow tests the complete wandfile workflow
func TestWandfileWorkflow(t *testing.T) {
	// Create isolated test home directory
	testHome := filepath.Join(".", "test-home-wandfile")

	// Clean up test directory before and after
	if err := os.RemoveAll(testHome); err != nil {
		t.Fatalf("Failed to clean test home: %v", err)
	}
	defer func() { _ = os.RemoveAll(testHome) }()

	if err := os.MkdirAll(testHome, 0755); err != nil { //nolint:gosec
		t.Fatalf("Failed to create test home: %v", err)
	}

	wandDir := filepath.Join(testHome, ".wand")
	formulasDir := filepath.Join(wandDir, "formulas")

	// Set up test environment with formulas
	if err := os.MkdirAll(formulasDir, 0755); err != nil { //nolint:gosec
		t.Fatalf("Failed to create formulas dir: %v", err)
	}

	// Copy formula files to test environment
	formulaFiles := []string{"nano.yaml", "zsh.yaml"}
	for _, formulaFile := range formulaFiles {
		sourceFormula := filepath.Join("..", "formulas", formulaFile)
		destFormula := filepath.Join(formulasDir, formulaFile)

		formulaData, err := os.ReadFile(sourceFormula) //nolint:gosec
		if err != nil {
			t.Logf("Warning: Failed to read %s: %v", formulaFile, err)
			continue
		}

		if err := os.WriteFile(destFormula, formulaData, 0644); err != nil { //nolint:gosec
			t.Fatalf("Failed to copy %s: %v", formulaFile, err)
		}
	}

	// Initialize all dependencies
	fs := domain_adapters.NewFileSystemAdapter()
	downloader := domain_adapters.NewDownloaderAdapter()
	extractor := domain_adapters.NewExtractorAdapter(fs)
	shellExecutor := domain_adapters.NewShellExecutorAdapter()

	// Initialize repositories
	registryRepo := domain_adapters.NewRegistryRepository(fs, wandDir)
	formulaRepo := domain_adapters.NewFormulaRepository(fs, formulasDir)
	wandrcRepo := domain_adapters.NewWandRCRepository(fs)
	wandfileRepo := domain_adapters.NewWandfileRepository(fs)
	dotfileRepo := domain_adapters.NewDotfileRepository(fs, wandDir)

	// Initialize external adapters
	githubClient := external_adapters.NewGitHubAdapter("")

	// Initialize domain services
	versionService := services.NewVersionService(githubClient, formulaRepo)
	_ = services.NewShimService(registryRepo, wandrcRepo, formulaRepo, fs, wandDir) // Not used in this test
	installerService := services.NewInstallerService(
		formulaRepo,
		registryRepo,
		downloader,
		extractor,
		fs,
		shellExecutor,
		versionService,
		wandDir,
		testHome,
	)
	wandfileService := services.NewWandfileService(
		wandfileRepo,
		registryRepo,
		installerService,
		versionService,
		dotfileRepo,
		fs,
		shellExecutor,
		testHome,
	)

	t.Run("CreateWandfile", func(t *testing.T) {
		wandfile := entities.NewWandfile()
		wandfile.AddCLI("nano", "8.7.0")
		wandfile.AddCLI("zsh", "5.9.0")

		wandfilePath := filepath.Join(testHome, "wandfile")
		if err := wandfileRepo.Save(wandfilePath, wandfile); err != nil {
			t.Fatalf("Failed to save wandfile: %v", err)
		}

		// Verify file exists
		if !wandfileRepo.Exists(wandfilePath) {
			t.Error("Wandfile not created")
		}

		t.Log("✓ Created wandfile successfully")
	})

	t.Run("LoadWandfile", func(t *testing.T) {
		wandfilePath := filepath.Join(testHome, "wandfile")
		wandfile, err := wandfileRepo.Load(wandfilePath)
		if err != nil {
			t.Fatalf("Failed to load wandfile: %v", err)
		}

		if len(wandfile.CLI) != 2 {
			t.Errorf("Expected 2 CLI packages, got %d", len(wandfile.CLI))
		}

		if wandfile.CLI[0].Name != "nano" || wandfile.CLI[0].Version != "8.7.0" {
			t.Errorf("Unexpected CLI package: %v", wandfile.CLI[0])
		}

		t.Logf("✓ Loaded wandfile with %d CLI packages", len(wandfile.CLI))
	})

	t.Run("CheckWandfileBeforeInstall", func(t *testing.T) {
		wandfilePath := filepath.Join(testHome, "wandfile")
		wandfile, err := wandfileRepo.Load(wandfilePath)
		if err != nil {
			t.Fatalf("Failed to load wandfile: %v", err)
		}

		missing, err := wandfileService.Check(wandfile)
		if err != nil {
			t.Fatalf("Failed to check wandfile: %v", err)
		}

		if len(missing) != 2 {
			t.Errorf("Expected 2 missing packages, got %d", len(missing))
		}

		t.Logf("✓ Check found %d missing packages (expected)", len(missing))
	})

	t.Run("DumpWandfile", func(t *testing.T) {
		// Create a test registry with some packages
		registry, err := registryRepo.Load()
		if err != nil {
			t.Fatalf("Failed to load registry: %v", err)
		}

		// Add some test packages to registry
		nanoVersion, _ := entities.NewVersion("8.7.0")
		nanoPkg := entities.NewPackage("nano", entities.PackageTypeCLI, nanoVersion)
		registry.AddPackage(nanoPkg)
		registry.SetGlobalVersion("nano", "8.7.0")

		zshVersion, _ := entities.NewVersion("5.9.0")
		zshPkg := entities.NewPackage("zsh", entities.PackageTypeCLI, zshVersion)
		registry.AddPackage(zshPkg)
		registry.SetGlobalVersion("zsh", "5.9.0")

		if err := registryRepo.Save(registry); err != nil {
			t.Fatalf("Failed to save registry: %v", err)
		}

		// Dump wandfile
		wandfile, err := wandfileService.Dump()
		if err != nil {
			t.Fatalf("Failed to dump wandfile: %v", err)
		}

		if len(wandfile.CLI) != 2 {
			t.Errorf("Expected 2 CLI packages in dump, got %d", len(wandfile.CLI))
		}

		// Save dumped wandfile
		dumpPath := filepath.Join(testHome, "wandfile-dump")
		if err := wandfileRepo.Save(dumpPath, wandfile); err != nil {
			t.Fatalf("Failed to save dumped wandfile: %v", err)
		}

		t.Logf("✓ Dumped wandfile with %d packages", len(wandfile.CLI))
	})
}
