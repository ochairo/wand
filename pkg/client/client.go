// Package client provides a high-level API for interacting with wand programmatically.
//
// This package is designed for third-party integrations, TUIs, web UIs, and other tools
// that need to interact with wand's functionality without using the CLI.
//
// Example usage:
//
//	import "github.com/ochairo/wand/pkg/client"
//
//	// Create a new client
//	c, err := client.New("")  // uses default ~/.wand
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Install a package
//	pkg, err := c.Install("jq", "1.7.1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// List installed packages
//	packages, err := c.ListPackages()
//	for _, pkg := range packages {
//	    fmt.Printf("%s@%s\n", pkg.Name, pkg.VersionString())
//	}
package client

import (
	"fmt"
	"os"
	"path/filepath"

	domainadapters "github.com/ochairo/wand/internal/domain-adapters"
	domainorchestrators "github.com/ochairo/wand/internal/domain-orchestrators"
	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
	"github.com/ochairo/wand/internal/domain/services"
	externaladapters "github.com/ochairo/wand/internal/external-adapters"
	"github.com/ochairo/wand/pkg/types"
)

// Client provides programmatic access to wand functionality.
type Client struct {
	wandDir string
	homeDir string

	// Internal services
	versionService   *services.VersionService
	installerService *services.InstallerService
	shimService      *services.ShimService
	wandfileService  *services.WandfileService

	// Orchestrators
	installOrchestrator *domainorchestrators.InstallOrchestrator

	// Repositories
	registryRepo interfaces.RegistryRepository
	formulaRepo  interfaces.FormulaRepository
	wandrcRepo   interfaces.WandRCRepository
	wandfileRepo interfaces.WandfileRepository
	dotfileRepo  interfaces.DotfileRepository
}

// New creates a new wand client.
// If wandDir is empty, uses the default ~/.wand directory.
func New(wandDir string) (*Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	if wandDir == "" {
		wandDir = filepath.Join(homeDir, ".wand")
	}

	formulasDir := filepath.Join(wandDir, "formulas")

	// Initialize adapters
	fs := domainadapters.NewFileSystemAdapter()
	downloader := domainadapters.NewDownloaderAdapter()
	extractor := domainadapters.NewExtractorAdapter(fs)
	shellExecutor := domainadapters.NewShellExecutorAdapter()

	// Initialize repositories
	registryRepo := domainadapters.NewRegistryRepository(fs, wandDir)
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulasDir)
	wandrcRepo := domainadapters.NewWandRCRepository(fs)
	wandfileRepo := domainadapters.NewWandfileRepository(fs)
	dotfileRepo := domainadapters.NewDotfileRepository(fs, wandDir)

	// Initialize external adapters
	githubClient := externaladapters.NewGitHubAdapter("")

	// Initialize services
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
		homeDir,
	)
	wandfileService := services.NewWandfileService(
		wandfileRepo,
		registryRepo,
		installerService,
		versionService,
		dotfileRepo,
		fs,
		shellExecutor,
		homeDir,
	)

	// Initialize orchestrators
	installOrchestrator := domainorchestrators.NewInstallOrchestrator(
		installerService,
		shimService,
		versionService,
		formulaRepo,
	)

	return &Client{
		wandDir:             wandDir,
		homeDir:             homeDir,
		versionService:      versionService,
		installerService:    installerService,
		shimService:         shimService,
		wandfileService:     wandfileService,
		installOrchestrator: installOrchestrator,
		registryRepo:        registryRepo,
		formulaRepo:         formulaRepo,
		wandrcRepo:          wandrcRepo,
		wandfileRepo:        wandfileRepo,
		dotfileRepo:         dotfileRepo,
	}, nil
}

// Install installs a package with the specified version.
// If version is empty or "latest", installs the latest available version.
func (c *Client) Install(packageName, version string) (*types.Package, error) {
	if version == "" {
		version = "latest"
	}

	err := c.installOrchestrator.InstallPackage(packageName, version)
	if err != nil {
		return nil, fmt.Errorf("installation failed: %w", err)
	}

	// Load registry to get the installed package
	registry, err := c.registryRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	entry, exists := registry.Packages[packageName]
	if !exists {
		return nil, fmt.Errorf("package not found in registry after installation")
	}

	pkg, ok := entry.Versions[version]
	if !ok {
		// If exact version not found, get the global version
		globalVer, _ := registry.GetGlobalVersion(packageName)
		pkg = entry.Versions[globalVer]
	}

	return convertPackage(pkg), nil
}

// Uninstall removes a specific version of a package.
func (c *Client) Uninstall(packageName, version string) error {
	return c.installOrchestrator.UninstallPackage(packageName, version)
}

// IsInstalled checks if a specific package version is installed.
func (c *Client) IsInstalled(packageName, version string) bool {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return false
	}

	entry, exists := registry.Packages[packageName]
	if !exists {
		return false
	}

	_, ok := entry.Versions[version]
	return ok
}

// ListInstalled returns all installed versions of a package.
func (c *Client) ListInstalled(packageName string) ([]*types.Package, error) {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	entry, exists := registry.Packages[packageName]
	if !exists {
		return []*types.Package{}, nil
	}

	var packages []*types.Package
	for _, pkg := range entry.Versions {
		packages = append(packages, convertPackage(pkg))
	}

	return packages, nil
}

// GetRegistry returns the current registry state.
func (c *Client) GetRegistry() (*types.Registry, error) {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	return convertRegistry(registry), nil
}

// SetGlobalVersion sets the active global version for a package.
func (c *Client) SetGlobalVersion(packageName, version string) error {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	registry.SetGlobalVersion(packageName, version)

	if err := c.registryRepo.Save(registry); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	return nil
}

// GetGlobalVersion returns the active global version for a package.
func (c *Client) GetGlobalVersion(packageName string) (string, error) {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load registry: %w", err)
	}

	version, ok := registry.GetGlobalVersion(packageName)
	if !ok {
		return "", fmt.Errorf("no global version set for %s", packageName)
	}

	return version, nil
}

// ListPackages returns all registered packages.
func (c *Client) ListPackages() ([]*types.PackageEntry, error) {
	registry, err := c.registryRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	var entries []*types.PackageEntry
	for _, entry := range registry.Packages {
		entries = append(entries, convertPackageEntry(entry))
	}

	return entries, nil
}

// GetFormula retrieves a formula by name.
func (c *Client) GetFormula(name string) (*types.Formula, error) {
	formula, err := c.formulaRepo.GetFormula(name)
	if err != nil {
		return nil, fmt.Errorf("formula not found: %w", err)
	}

	return convertFormula(formula), nil
}

// ListFormulas returns all available formulas.
func (c *Client) ListFormulas() ([]*types.Formula, error) {
	formulas, err := c.formulaRepo.ListFormulas()
	if err != nil {
		return nil, fmt.Errorf("failed to list formulas: %w", err)
	}

	var result []*types.Formula
	for _, formula := range formulas {
		result = append(result, convertFormula(formula))
	}

	return result, nil
}

// SearchFormulas searches formulas by name or tag.
func (c *Client) SearchFormulas(query string) ([]*types.Formula, error) {
	allFormulas, err := c.ListFormulas()
	if err != nil {
		return nil, err
	}

	var matches []*types.Formula
	for _, formula := range allFormulas {
		// Simple search: check if query is in name or tags
		if contains(formula.Name, query) {
			matches = append(matches, formula)
			continue
		}
		for _, tag := range formula.Tags {
			if contains(tag, query) {
				matches = append(matches, formula)
				break
			}
		}
	}

	return matches, nil
}

// ListAvailableVersions returns all available versions for a package.
func (c *Client) ListAvailableVersions(packageName string) ([]*types.Version, error) {
	versions, err := c.versionService.ListAvailableVersions(packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	var result []*types.Version
	for _, v := range versions {
		result = append(result, convertVersion(v))
	}

	return result, nil
}

// ResolveVersion resolves a version constraint to a specific version.
func (c *Client) ResolveVersion(packageName, constraint string) (*types.Version, error) {
	version, err := c.versionService.ResolveVersion(packageName, constraint)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve version: %w", err)
	}

	return convertVersion(version), nil
}

// Helper conversion functions
func convertPackage(pkg *entities.Package) *types.Package {
	return &types.Package{
		Name:        pkg.Name,
		Type:        types.PackageType(pkg.Type),
		Version:     convertVersion(pkg.Version),
		InstalledAt: pkg.InstalledAt,
		BinPath:     pkg.BinPath,
		InstallPath: pkg.InstallPath,
		IsGlobal:    pkg.IsGlobal,
	}
}

func convertVersion(v *entities.Version) *types.Version {
	if v == nil {
		return nil
	}
	return &types.Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
		Pre:   v.Pre,
		Build: v.Build,
	}
}

func convertRegistry(r *entities.Registry) *types.Registry {
	registry := &types.Registry{
		Packages:       make(map[string]*types.PackageEntry),
		GlobalVersions: r.GlobalVersions,
		UpdatedAt:      r.UpdatedAt,
	}

	for name, entry := range r.Packages {
		registry.Packages[name] = convertPackageEntry(entry)
	}

	return registry
}

func convertPackageEntry(e *entities.PackageEntry) *types.PackageEntry {
	entry := &types.PackageEntry{
		Name:     e.Name,
		Type:     types.PackageType(e.Type),
		Versions: make(map[string]*types.Package),
	}

	for version, pkg := range e.Versions {
		entry.Versions[version] = convertPackage(pkg)
	}

	return entry
}

func convertFormula(f *entities.Formula) *types.Formula {
	return &types.Formula{
		Name:           f.Name,
		Type:           types.PackageType(f.Type),
		Description:    f.Description,
		Homepage:       f.Homepage,
		Repository:     f.Repository,
		License:        f.License,
		Tags:           f.Tags,
		VersionPattern: f.VersionPattern,
		Binaries:       f.Binaries,
		BinPath:        f.BinPath,
		AppName:        f.AppName,
		AppPath:        f.AppPath,
		MinVersion:     f.MinVersion,
		MaxVersion:     f.MaxVersion,
		Dependencies:   f.Dependencies,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
