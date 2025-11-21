// Package services provides domain business logic services.
package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	errs "github.com/ochairo/wand/internal/domain/errors"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// InstallerService handles package installation logic
type InstallerService struct {
	formulaRepo   interfaces.FormulaRepository
	registryRepo  interfaces.RegistryRepository
	downloader    interfaces.Downloader
	extractor     interfaces.Extractor
	fs            interfaces.FileSystem
	shellExecutor interfaces.ShellExecutor
	versionSvc    *VersionService
	wandDir       string
	homeDir       string
}

// NewInstallerService creates a new installer service
func NewInstallerService(
	formulaRepo interfaces.FormulaRepository,
	registryRepo interfaces.RegistryRepository,
	downloader interfaces.Downloader,
	extractor interfaces.Extractor,
	fs interfaces.FileSystem,
	shellExecutor interfaces.ShellExecutor,
	versionSvc *VersionService,
	wandDir string,
	homeDir string,
) *InstallerService {
	return &InstallerService{
		formulaRepo:   formulaRepo,
		registryRepo:  registryRepo,
		downloader:    downloader,
		extractor:     extractor,
		fs:            fs,
		shellExecutor: shellExecutor,
		versionSvc:    versionSvc,
		wandDir:       wandDir,
		homeDir:       homeDir,
	}
}

// InstallPackage installs a package with a specific version
func (s *InstallerService) InstallPackage(packageName, versionStr string) error {
	// Get formula
	formula, err := s.formulaRepo.GetFormula(packageName)
	if err != nil {
		return errs.New(errs.ErrPackageNotFound, fmt.Sprintf("Formula not found for package %q", packageName))
	}

	// Resolve version
	version, err := s.versionSvc.ResolveVersion(packageName, versionStr)
	if err != nil {
		return err // propagate from VersionService
	}

	// Check if already installed
	registry, err := s.registryRepo.Load()
	if err != nil && !s.registryRepo.Exists() {
		registry = entities.NewRegistry()
	} else if err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry", err)
	}

	if _, exists := registry.GetPackage(packageName, version.String()); exists {
		return errs.NewWithDetails(errs.ErrPackageInstalled, "Package already installed", fmt.Sprintf("package: %q, version: %q", packageName, version.String()))
	}

	// Get platform config
	platform := entities.CurrentPlatform()
	platformConfig := formula.GetPlatformConfigFor(platform)
	if platformConfig == nil {
		return errs.NewWithDetails(errs.ErrArchNotSupported, "Package not available for platform", fmt.Sprintf("package: %q, os: %s, arch: %s", packageName, platform.OS, platform.Arch))
	}

	// Build download URL
	downloadURL := buildDownloadURL(platformConfig.DownloadURL, version, platform)

	// Create temp directory for download
	tmpDir := filepath.Join(s.wandDir, "tmp", packageName+"-"+version.String())
	if err := s.fs.MkdirAll(tmpDir, 0755); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to create temp directory", err)
	}
	defer func() { _ = s.fs.RemoveAll(tmpDir) }()

	// Preserve file extension from URL for proper extraction detection
	ext := filepath.Ext(downloadURL)
	if ext == ".gz" && strings.HasSuffix(downloadURL, ".tar.gz") {
		ext = ".tar.gz"
	}
	downloadPath := filepath.Join(tmpDir, "package"+ext)
	if err := s.downloader.Download(downloadURL, downloadPath); err != nil {
		return errs.Wrap(errs.ErrDownloadFailed, fmt.Sprintf("Failed to download %s@%s", packageName, version.String()), err)
	}

	// Verify checksum if checksum URL is provided
	if platformConfig.ChecksumURL != "" {
		checksumURL := buildDownloadURL(platformConfig.ChecksumURL, version, platform)
		if err := s.downloader.VerifyChecksum(downloadPath, checksumURL); err != nil {
			return errs.Wrap(errs.ErrChecksumMismatch, fmt.Sprintf("Checksum verification failed for %q", packageName), err)
		}
	}

	// Install based on package type
	switch formula.Type {
	case entities.PackageTypeCLI:
		return s.installCLI(formula, version, downloadPath, platformConfig)
	case entities.PackageTypeGUI:
		return s.installGUI(formula, version, downloadPath, platformConfig, platform)
	default:
		return errs.New(errs.ErrInstallationFailed, fmt.Sprintf("Unsupported package type: %s", formula.Type))
	}
}

// installCLI installs a CLI package
func (s *InstallerService) installCLI(
	formula *entities.Formula,
	version *entities.Version,
	downloadPath string,
	config *entities.PlatformConfig,
) error {
	// Create version directory
	installDir := filepath.Join(s.wandDir, "packages", formula.Name, version.String())
	if err := s.fs.MkdirAll(installDir, 0755); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to create install directory for %s@%s", formula.Name, version.String()), err)
	}

	// Extract if archive
	if isArchive(downloadPath) {
		if err := s.extractor.Extract(downloadPath, installDir); err != nil {
			return errs.Wrap(errs.ErrExtractionFailed, fmt.Sprintf("Failed to extract package %s@%s", formula.Name, version.String()), err)
		}
	} else {
		// Single binary
		binDir := filepath.Join(installDir, "bin")
		if err := s.fs.MkdirAll(binDir, 0755); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to create bin directory", err)
		}

		binPath := filepath.Join(binDir, formula.Binaries[0])
		data, err := s.fs.ReadFile(downloadPath)
		if err != nil {
			return errs.Wrap(errs.ErrFileNotFound, "Failed to read binary", err)
		}

		if err := s.fs.WriteFile(binPath, data, 0755); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to write binary", err)
		}
	}

	// Build from source if needed
	if config.RequiresBuild {
		if err := s.buildFromSource(installDir, config.BuildCommands); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, "Failed to build from source", err)
		}
	}

	// Run post-install hooks
	if formula.PostInstall != nil {
		if err := s.runPostInstall(installDir, formula.PostInstall); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, "Post-install hook failed", err)
		}
	}

	// Update registry
	if err := s.addToRegistry(formula.Name, version.String(), entities.PackageTypeCLI, installDir); err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to update registry", err)
	}

	return nil
}

// installGUI installs a GUI application
func (s *InstallerService) installGUI(
	formula *entities.Formula,
	version *entities.Version,
	downloadPath string,
	config *entities.PlatformConfig,
	platform *entities.Platform,
) error {
	appsDir := filepath.Join(s.wandDir, "apps", formula.Name)
	if err := s.fs.MkdirAll(appsDir, 0755); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to create apps directory", err)
	}

	// Extract application
	if err := s.extractor.Extract(downloadPath, appsDir); err != nil {
		return errs.Wrap(errs.ErrExtractionFailed, fmt.Sprintf("Failed to extract application %s@%s", formula.Name, version.String()), err)
	}

	if platform.IsDarwin() {
		// macOS: Symlink .app bundle to ~/Applications
		homeApps := filepath.Join(s.homeDir, "Applications")
		_ = s.fs.MkdirAll(homeApps, 0755)

		appPath := filepath.Join(appsDir, formula.AppName)
		symlinkPath := filepath.Join(homeApps, formula.AppName)

		if err := s.fs.Symlink(appPath, symlinkPath); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to create symlink for %s", formula.AppName), err)
		}
	} else if config.DesktopFile != "" {
		// Linux: Install .desktop file if provided
		desktopSrc := filepath.Join(appsDir, config.DesktopFile)
		desktopDir := filepath.Join(s.homeDir, ".local/share/applications")
		desktopDest := filepath.Join(desktopDir, formula.Name+".desktop")

		if err := s.fs.MkdirAll(desktopDir, 0755); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to create desktop directory", err)
		}

		data, err := s.fs.ReadFile(desktopSrc)
		if err != nil {
			return errs.Wrap(errs.ErrFileNotFound, "Failed to read desktop file", err)
		}

		if err := s.fs.WriteFile(desktopDest, data, 0644); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, "Failed to install desktop file", err)
		}
	}

	// Update registry
	if err := s.addToRegistry(formula.Name, version.String(), entities.PackageTypeGUI, appsDir); err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to update registry", err)
	}

	return nil
}

// buildFromSource runs build commands in the install directory
func (s *InstallerService) buildFromSource(dir string, commands []string) error {
	for _, cmd := range commands {
		if _, err := s.shellExecutor.ExecuteInDir(dir, cmd); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Build command failed: %s", cmd), err)
		}
	}
	return nil
}

// runPostInstall runs post-install hooks
func (s *InstallerService) runPostInstall(installDir string, hook *entities.PostInstallHook) error {
	binPath := filepath.Join(installDir, "bin")

	for _, cmd := range hook.Commands {
		// Replace {bin_path} placeholder
		cmdWithPath := strings.ReplaceAll(cmd, "{bin_path}", binPath)

		if _, err := s.shellExecutor.ExecuteWithEnv(hook.Env, cmdWithPath); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Post-install command failed: %s", cmd), err)
		}
	}
	return nil
}

// addToRegistry adds a package to the registry
func (s *InstallerService) addToRegistry(packageName, versionStr string, pkgType entities.PackageType, installDir string) error {
	registry, err := s.registryRepo.Load()
	if err != nil && !s.registryRepo.Exists() {
		registry = entities.NewRegistry()
	} else if err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry", err)
	}

	// Create version
	version, err := entities.NewVersion(versionStr)
	if err != nil {
		return errs.New(errs.ErrInvalidVersion, fmt.Sprintf("Invalid version: %q", versionStr))
	}

	// Create package
	pkg := entities.NewPackage(packageName, pkgType, version)
	pkg.InstallPath = installDir

	// Find binaries - they may be in bin/ subdirectory or at root level
	binDir := filepath.Join(installDir, "bin")
	if !s.fs.Exists(binDir) {
		// Binaries are at root level, use install dir
		pkg.BinPath = installDir
	} else {
		pkg.BinPath = binDir
	}

	pkg.IsGlobal = true

	// Add to registry
	registry.AddPackage(pkg)
	registry.SetGlobalVersion(packageName, versionStr)

	return s.registryRepo.Save(registry)
}

// UninstallPackage removes a specific version of a package
func (s *InstallerService) UninstallPackage(packageName, version string) error {
	registry, err := s.registryRepo.Load()
	if err != nil {
		return errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry for uninstall", err)
	}

	// If no version specified, uninstall all versions
	if version == "" {
		return s.uninstallAllVersions(registry, packageName)
	}

	pkg, exists := registry.GetPackage(packageName, version)
	if !exists {
		return errs.NewWithDetails(errs.ErrPackageNotInstalled, "Package not installed", fmt.Sprintf("package: %q, version: %q", packageName, version))
	}

	// Remove installation directory
	if err := s.fs.RemoveAll(pkg.InstallPath); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to remove installation", err)
	}

	// Update registry
	registry.RemovePackage(packageName, version)

	return s.registryRepo.Save(registry)
}

// uninstallAllVersions removes all installed versions of a package
func (s *InstallerService) uninstallAllVersions(registry *entities.Registry, packageName string) error {
	entry, exists := registry.Packages[packageName]
	if !exists {
		return errs.NewWithDetails(errs.ErrPackageNotInstalled, "Package not installed", fmt.Sprintf("package: %q", packageName))
	}

	if len(entry.Versions) == 0 {
		return errs.NewWithDetails(errs.ErrPackageNotInstalled, "No versions installed", fmt.Sprintf("package: %q", packageName))
	}

	// Remove all version directories
	for version, pkg := range entry.Versions {
		if err := s.fs.RemoveAll(pkg.InstallPath); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to remove %s@%s", packageName, version), err)
		}
	}

	// Remove entire package entry from registry
	delete(registry.Packages, packageName)
	delete(registry.GlobalVersions, packageName)

	return s.registryRepo.Save(registry)
}

// buildDownloadURL replaces placeholders in download URL template
func buildDownloadURL(template string, version *entities.Version, platform *entities.Platform) string {
	url := template

	replacements := map[string]string{
		"{version}":       version.ShortString(), // Use short version (e.g., 8.7 instead of 8.7.0)
		"{version_major}": fmt.Sprintf("%d", version.Major),
		"{version_minor}": fmt.Sprintf("%d.%d", version.Major, version.Minor),
		"{platform}":      platform.OS,
		"{os}":            platform.OS,
		"{arch}":          platform.Arch,
	}

	for placeholder, value := range replacements {
		url = strings.ReplaceAll(url, placeholder, value)
	}

	return url
}

// isArchive checks if a file is an archive based on extension
func isArchive(path string) bool {
	archiveExts := []string{".tar", ".gz", ".tgz", ".zip", ".bz2", ".xz", ".tar.gz", ".tar.bz2", ".tar.xz"}

	for _, ext := range archiveExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
