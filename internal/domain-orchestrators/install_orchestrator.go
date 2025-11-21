package domainorchestrators

import (
	"fmt"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
	"github.com/ochairo/wand/internal/domain/services"
)

// InstallOrchestrator orchestrates the package installation process
type InstallOrchestrator struct {
	installerSvc *services.InstallerService
	shimSvc      *services.ShimService
	versionSvc   *services.VersionService
	formulaRepo  interfaces.FormulaRepository
}

// NewInstallOrchestrator creates a new InstallOrchestrator
func NewInstallOrchestrator(
	installerSvc *services.InstallerService,
	shimSvc *services.ShimService,
	versionSvc *services.VersionService,
	formulaRepo interfaces.FormulaRepository,
) *InstallOrchestrator {
	return &InstallOrchestrator{
		installerSvc: installerSvc,
		shimSvc:      shimSvc,
		versionSvc:   versionSvc,
		formulaRepo:  formulaRepo,
	}
}

// InstallPackageOptions contains installation options
type InstallPackageOptions struct {
	Global bool // Install to global location (/usr/local/bin)
	Force  bool // Overwrite existing installation
}

// InstallPackageWithOptions installs a package with the specified options and creates shims for all binaries.
// InstallPackage installs a package and creates shims
func (o *InstallOrchestrator) InstallPackageWithOptions(packageName, versionStr string, opts InstallPackageOptions) error {
	// If force flag is set, uninstall existing version first
	if opts.Force {
		// Try to uninstall, but don't fail if it doesn't exist
		_ = o.UninstallPackage(packageName, "*")
	}

	// Install the package
	if err := o.installerSvc.InstallPackage(packageName, versionStr); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Load formula to get actual binary names
	formula, err := o.formulaRepo.GetFormula(packageName)
	if err != nil {
		// Fallback: assume binary name matches package name
		formula = nil
	}

	// Get binaries from formula or use package name as fallback
	var binaries []string
	if formula != nil && len(formula.Binaries) > 0 {
		binaries = formula.Binaries
	} else {
		binaries = []string{packageName}
	}

	// Create shims for CLI packages
	if err := o.shimSvc.CreateShims(packageName, binaries); err != nil {
		return fmt.Errorf("failed to create shims: %w", err)
	}

	return nil
}

// InstallPackage installs a package and creates shims (backward compatible)
func (o *InstallOrchestrator) InstallPackage(packageName, versionStr string) error {
	return o.InstallPackageWithOptions(packageName, versionStr, InstallPackageOptions{})
}

// UninstallPackage removes a package and its shims
func (o *InstallOrchestrator) UninstallPackage(packageName, version string) error {
	// Load formula to get actual binary names
	formula, err := o.formulaRepo.GetFormula(packageName)
	if err != nil {
		// Fallback: assume binary name matches package name
		formula = nil
	}

	// Get binaries from formula or use package name as fallback
	var binaries []string
	if formula != nil && len(formula.Binaries) > 0 {
		binaries = formula.Binaries
	} else {
		binaries = []string{packageName}
	}

	// Remove shims first
	if err := o.shimSvc.RemoveShims(binaries); err != nil {
		// Log warning but continue
		fmt.Printf("Warning: failed to remove shims: %v\n", err)
	}

	// Uninstall the package
	if err := o.installerSvc.UninstallPackage(packageName, version); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}

	return nil
}

// ListAvailableVersions lists all available versions from the repository
func (o *InstallOrchestrator) ListAvailableVersions(packageName string) ([]*entities.Version, error) {
	return o.versionSvc.ListAvailableVersions(packageName)
}
