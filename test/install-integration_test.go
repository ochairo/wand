package test //nolint:dupl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	domain_adapters "github.com/ochairo/wand/internal/domain-adapters"
	domain_orchestrators "github.com/ochairo/wand/internal/domain-orchestrators"
	"github.com/ochairo/wand/internal/domain/services"
	external_adapters "github.com/ochairo/wand/internal/external-adapters"
)

// TestWandSetup verifies that the wand system can be initialized correctly in isolation
func TestWandSetup(t *testing.T) {
	// Create isolated test home directory - NEVER touches real home directory
	testHome := filepath.Join(".", "test-home-setup")

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

	// Copy all formula files to test environment
	formulaFiles := []string{"nano.yaml", "zsh.yaml", "make.yaml", "zlib.yaml", "microsoft-edge.yaml"}
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

	// Initialize all dependencies exactly like main.go but with test paths
	fs := domain_adapters.NewFileSystemAdapter()
	downloader := domain_adapters.NewDownloaderAdapter()
	extractor := domain_adapters.NewExtractorAdapter(fs)
	shellExecutor := domain_adapters.NewShellExecutorAdapter()

	// Initialize repositories with test paths
	registryRepo := domain_adapters.NewRegistryRepository(fs, wandDir)
	formulaRepo := domain_adapters.NewFormulaRepository(fs, formulasDir)
	wandrcRepo := domain_adapters.NewWandRCRepository(fs)

	// Initialize external adapters
	githubClient := external_adapters.NewGitHubAdapter("")

	// Initialize domain services
	versionService := services.NewVersionService(githubClient, formulaRepo)
	shimService := services.NewShimService(registryRepo, wandrcRepo, formulaRepo, fs, wandDir)
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

	// Initialize orchestrator
	installOrchestrator := domain_orchestrators.NewInstallOrchestrator(
		installerService,
		shimService,
		versionService,
		formulaRepo,
	)

	t.Run("LoadFormulas", func(t *testing.T) {
		formula, err := formulaRepo.GetFormula("nano")
		if err != nil {
			t.Fatalf("Failed to load nano formula: %v", err)
		}

		if formula.Name != "nano" {
			t.Errorf("Expected formula name 'nano', got '%s'", formula.Name)
		}

		if formula.Type != "cli" {
			t.Errorf("Expected type 'cli', got '%s'", formula.Type)
		}

		if len(formula.Binaries) == 0 {
			t.Error("Expected at least one binary defined")
		}

		t.Logf("✓ Formula loaded successfully: %s - %s", formula.Name, formula.Description)
	})

	t.Run("RegistryOperations", func(t *testing.T) {
		// Load empty registry
		registry, err := registryRepo.Load()
		if err != nil {
			t.Fatalf("Failed to load registry: %v", err)
		}

		if len(registry.Packages) != 0 {
			t.Error("Expected empty registry initially")
		}

		// Save registry
		err = registryRepo.Save(registry)
		if err != nil {
			t.Fatalf("Failed to save registry: %v", err)
		}

		// Verify file exists
		registryPath := filepath.Join(wandDir, "registry.json")
		if _, err := os.Stat(registryPath); os.IsNotExist(err) {
			t.Error("Registry file not created")
		}

		t.Logf("✓ Registry operations work correctly")
	})

	t.Run("CheckGitHubAPI", func(t *testing.T) {
		// Test GitHub API connectivity (without installing)
		releases, err := githubClient.ListReleases("ochairo", "potions")
		if err != nil {
			t.Logf("Warning: GitHub API not accessible: %v", err)
			t.Skip("Skipping GitHub API test")
		}

		t.Logf("✓ GitHub API accessible, found %d releases", len(releases))

		// Log all release tags to see if nano exists
		nanoFound := false
		for _, release := range releases {
			t.Logf("  Release tag: %s", release.TagName)
			if strings.HasPrefix(release.TagName, "nano-") {
				nanoFound = true
			}
		}

		if !nanoFound {
			t.Log("Note: No nano releases found in repository")
		}
	})

	t.Run("AvailableVersions", func(t *testing.T) {
		// Try to list available versions for nano
		versions, err := installOrchestrator.ListAvailableVersions("nano")
		if err != nil {
			t.Logf("Note: Could not fetch versions for nano: %v", err)
			t.Log("This is expected if ochairo/potions doesn't have nano releases yet")
			return
		}

		t.Logf("✓ Found %d versions available for nano", len(versions))
		for _, v := range versions {
			t.Logf("  - %s", v.String())
		}
	})

	t.Run("InstallNanoLatest", func(t *testing.T) {
		t.Skip("Integration test - requires real GitHub API and formula versions")
	})

	t.Run("UninstallNano", func(t *testing.T) {
		// Get installed version from registry
		registry, err := registryRepo.Load()
		if err != nil {
			t.Fatalf("Failed to load registry: %v", err)
		}

		pkg, found := registry.Packages["nano"]
		if !found {
			t.Skip("nano not installed, skipping uninstall test")
		}

		var version string
		for v := range pkg.Versions {
			version = v
			break
		}

		err = installOrchestrator.UninstallPackage("nano", version)
		if err != nil {
			t.Fatalf("Failed to uninstall nano@%s: %v", version, err)
		}

		// Verify package removed from test-home
		nanoDir := filepath.Join(wandDir, "packages", "nano", version)
		if _, err := os.Stat(nanoDir); !os.IsNotExist(err) {
			t.Error("nano directory still exists after uninstall")
		}

		// Verify shim removed
		shimPath := filepath.Join(wandDir, "shims", "nano")
		if _, err := os.Stat(shimPath); !os.IsNotExist(err) {
			t.Error("nano shim still exists after uninstalling all versions")
		}

		t.Logf("✓ Successfully uninstalled nano@%s from test environment", version)
	})
}
