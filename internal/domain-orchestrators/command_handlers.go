// Package domainorchestrators provides orchestrator implementations for domain workflows.
package domainorchestrators

import (
	"fmt"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// InstallCommandHandler handles the install command
type InstallCommandHandler struct {
	installOrchestrator *InstallOrchestrator
	registryRepo        interfaces.RegistryRepository
	wandrcRepo          interfaces.WandRCRepository
}

// NewInstallCommandHandler creates a new install command handler
func NewInstallCommandHandler(
	installOrchestrator *InstallOrchestrator,
	registryRepo interfaces.RegistryRepository,
	wandrcRepo interfaces.WandRCRepository,
) *InstallCommandHandler {
	return &InstallCommandHandler{
		installOrchestrator: installOrchestrator,
		registryRepo:        registryRepo,
		wandrcRepo:          wandrcRepo,
	}
}

// Handle executes the install command
func (h *InstallCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("package name required")
	}

	// Parse package[@version]
	packageSpec := args[0]
	parts := strings.Split(packageSpec, "@")
	packageName := parts[0]
	versionStr := "latest"
	if len(parts) > 1 {
		versionStr = parts[1]
	}

	// Get flags
	globalFlag, err := ctx.GetBoolFlag("global")
	if err != nil {
		globalFlag = false // default to false if flag not found
	}

	forceFlag, err := ctx.GetBoolFlag("force")
	if err != nil {
		forceFlag = false // default to false if flag not found
	}

	ctx.Printf("Installing %s@%s...\n", packageName, versionStr)

	// Install with flags
	opts := InstallPackageOptions{
		Global: globalFlag,
		Force:  forceFlag,
	}
	if err := h.installOrchestrator.InstallPackageWithOptions(packageName, versionStr, opts); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	ctx.Printf("✓ Successfully installed %s@%s\n", packageName, versionStr)

	return nil
}

// ListCommandHandler handles the list command
type ListCommandHandler struct {
	registryRepo interfaces.RegistryRepository
	wandrcRepo   interfaces.WandRCRepository
}

// NewListCommandHandler creates a new list command handler
func NewListCommandHandler(
	registryRepo interfaces.RegistryRepository,
	wandrcRepo interfaces.WandRCRepository,
) *ListCommandHandler {
	return &ListCommandHandler{
		registryRepo: registryRepo,
		wandrcRepo:   wandrcRepo,
	}
}

// Handle executes the list command
func (h *ListCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()

	// Load registry
	registry, err := h.registryRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// If package name provided, show versions for that package
	if len(args) > 0 {
		packageName := args[0]
		entry, exists := registry.Packages[packageName]
		if !exists {
			ctx.Printf("Package '%s' is not installed\n", packageName)
			return nil
		}

		ctx.Printf("Installed versions of %s:\n", packageName)
		for version, pkg := range entry.Versions {
			marker := " "
			if version == registry.GlobalVersions[packageName] {
				marker = "*"
			}
			ctx.Printf("  %s %s (installed at %s)\n", marker, version, pkg.InstalledAt.Format("2006-01-02"))
		}
		return nil
	}

	// Show all installed packages
	if len(registry.Packages) == 0 {
		ctx.Printf("No packages installed\n")
		return nil
	}

	ctx.Printf("Installed packages:\n")
	for name, entry := range registry.Packages {
		globalVer := registry.GlobalVersions[name]
		ctx.Printf("\n%s:\n", name)
		for version, pkg := range entry.Versions {
			marker := " "
			if version == globalVer {
				marker = "*"
			}
			ctx.Printf("  %s %s (%s)\n", marker, version, pkg.Type)
		}
	}

	return nil
}

// UseCommandHandler handles the use command
type UseCommandHandler struct {
	registryRepo interfaces.RegistryRepository
	wandrcRepo   interfaces.WandRCRepository
}

// NewUseCommandHandler creates a new use command handler
func NewUseCommandHandler(
	registryRepo interfaces.RegistryRepository,
	wandrcRepo interfaces.WandRCRepository,
) *UseCommandHandler {
	return &UseCommandHandler{
		registryRepo: registryRepo,
		wandrcRepo:   wandrcRepo,
	}
}

// Handle executes the use command
func (h *UseCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) < 2 {
		return fmt.Errorf("usage: wand use <package> <version>")
	}

	packageName := args[0]
	versionStr := args[1]

	global, err := ctx.GetBoolFlag("global")
	if err != nil {
		return fmt.Errorf("failed to get global flag: %w", err)
	}

	// Parse version
	version, err := entities.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Verify package version is installed
	registry, err := h.registryRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	entry, exists := registry.Packages[packageName]
	if !exists {
		return fmt.Errorf("package '%s' is not installed", packageName)
	}

	if _, exists := entry.Versions[version.String()]; !exists {
		return fmt.Errorf("version %s of package '%s' is not installed", version, packageName)
	}

	if global {
		// Set global version
		registry.GlobalVersions[packageName] = version.String()
		if err := h.registryRepo.Save(registry); err != nil {
			return fmt.Errorf("failed to save registry: %w", err)
		}
		ctx.Printf("✓ Set global version of %s to %s\n", packageName, version)
	} else {
		// Set project version in .wandrc
		wandrc, err := h.wandrcRepo.Load(".")
		if err != nil {
			// Create new .wandrc if it doesn't exist
			wandrc = &entities.WandRC{
				Versions: make(map[string]string),
			}
		}

		wandrc.Versions[packageName] = version.String()
		if err := h.wandrcRepo.Save(".", wandrc); err != nil {
			return fmt.Errorf("failed to save .wandrc: %w", err)
		}
		ctx.Printf("✓ Set project version of %s to %s\n", packageName, version)
	}

	return nil
}

// UninstallCommandHandler handles the uninstall command
type UninstallCommandHandler struct {
	uninstallOrchestrator *InstallOrchestrator
}

// NewUninstallCommandHandler creates a new uninstall command handler
func NewUninstallCommandHandler(uninstallOrchestrator *InstallOrchestrator) *UninstallCommandHandler {
	return &UninstallCommandHandler{
		uninstallOrchestrator: uninstallOrchestrator,
	}
}

// Handle executes the uninstall command
func (h *UninstallCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("package name required")
	}

	// Parse package[@version]
	packageSpec := args[0]
	parts := strings.Split(packageSpec, "@")
	packageName := parts[0]
	versionStr := ""
	if len(parts) > 1 {
		versionStr = parts[1]
	}

	if versionStr != "" {
		ctx.Printf("Uninstalling %s@%s...\n", packageName, versionStr)
	} else {
		ctx.Printf("Uninstalling all versions of %s...\n", packageName)
	}

	// Uninstall the package
	if err := h.uninstallOrchestrator.UninstallPackage(packageName, versionStr); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}

	if versionStr != "" {
		ctx.Printf("✓ Successfully uninstalled %s@%s\n", packageName, versionStr)
	} else {
		ctx.Printf("✓ Successfully uninstalled all versions of %s\n", packageName)
	}

	return nil
}

// InitCommandHandler handles the init command
type InitCommandHandler struct {
	wandrcRepo interfaces.WandRCRepository
}

// NewInitCommandHandler creates a new init command handler
func NewInitCommandHandler(wandrcRepo interfaces.WandRCRepository) *InitCommandHandler {
	return &InitCommandHandler{
		wandrcRepo: wandrcRepo,
	}
}

// Handle executes the init command
func (h *InitCommandHandler) Handle(ctx interfaces.CommandContext) error {
	// Check if .wandrc already exists
	if h.wandrcRepo.Exists(".") {
		return fmt.Errorf(".wandrc already exists in current directory")
	}

	// Create new .wandrc
	wandrc := entities.NewWandRC()
	if err := h.wandrcRepo.Save(".", wandrc); err != nil {
		return fmt.Errorf("failed to create .wandrc: %w", err)
	}

	ctx.Printf("✓ Created .wandrc in current directory\n")
	ctx.Printf("\nUse 'wand add <package>@<version>' to pin package versions for this project.\n")

	return nil
}

// AddCommandHandler handles the add command
type AddCommandHandler struct {
	wandrcRepo   interfaces.WandRCRepository
	registryRepo interfaces.RegistryRepository
}

// NewAddCommandHandler creates a new add command handler
func NewAddCommandHandler(
	wandrcRepo interfaces.WandRCRepository,
	registryRepo interfaces.RegistryRepository,
) *AddCommandHandler {
	return &AddCommandHandler{
		wandrcRepo:   wandrcRepo,
		registryRepo: registryRepo,
	}
}

// Handle executes the add command
func (h *AddCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("package specification required (e.g., nano@8.7.0)")
	}

	// Parse package[@version]
	packageSpec := args[0]
	parts := strings.Split(packageSpec, "@")
	if len(parts) != 2 {
		return fmt.Errorf("version required for add command (use package@version format)")
	}

	packageName := parts[0]
	versionStr := parts[1]

	// Parse and validate version
	version, err := entities.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Verify package version is installed
	registry, err := h.registryRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	entry, exists := registry.Packages[packageName]
	if !exists {
		return fmt.Errorf("package '%s' is not installed. Install it first with: wand install %s@%s", packageName, packageName, version)
	}

	if _, exists := entry.Versions[version.String()]; !exists {
		return fmt.Errorf("version %s of '%s' is not installed. Install it first with: wand install %s@%s", version, packageName, packageName, version)
	}

	// Load or create .wandrc
	var wandrc *entities.WandRC
	if h.wandrcRepo.Exists(".") {
		wandrc, err = h.wandrcRepo.Load(".")
		if err != nil {
			return fmt.Errorf("failed to load .wandrc: %w", err)
		}
	} else {
		wandrc = entities.NewWandRC()
	}

	// Add package version
	wandrc.SetVersion(packageName, version.String())

	// Save .wandrc
	if err := h.wandrcRepo.Save(".", wandrc); err != nil {
		return fmt.Errorf("failed to save .wandrc: %w", err)
	}

	ctx.Printf("✓ Added %s@%s to .wandrc\n", packageName, version)

	return nil
}

// RemoveCommandHandler handles the remove command
type RemoveCommandHandler struct {
	wandrcRepo interfaces.WandRCRepository
}

// NewRemoveCommandHandler creates a new remove command handler
func NewRemoveCommandHandler(wandrcRepo interfaces.WandRCRepository) *RemoveCommandHandler {
	return &RemoveCommandHandler{
		wandrcRepo: wandrcRepo,
	}
}

// Handle executes the remove command
func (h *RemoveCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("package name required")
	}

	packageName := args[0]

	// Check if .wandrc exists
	if !h.wandrcRepo.Exists(".") {
		return fmt.Errorf(".wandrc not found in current directory")
	}

	// Load .wandrc
	wandrc, err := h.wandrcRepo.Load(".")
	if err != nil {
		return fmt.Errorf("failed to load .wandrc: %w", err)
	}

	// Check if package exists in .wandrc
	if !wandrc.HasVersion(packageName) {
		return fmt.Errorf("package '%s' not found in .wandrc", packageName)
	}

	// Remove package
	wandrc.RemoveVersion(packageName)

	// Save .wandrc
	if err := h.wandrcRepo.Save(".", wandrc); err != nil {
		return fmt.Errorf("failed to save .wandrc: %w", err)
	}

	ctx.Printf("✓ Removed %s from .wandrc\n", packageName)

	return nil
}

// WandfileInstallCommandHandler handles the wandfile install command
type WandfileInstallCommandHandler struct {
	wandfileRepo interfaces.WandfileRepository
	wandfileSvc  interfaces.WandfileManager
}

// NewWandfileInstallCommandHandler creates a new wandfile install command handler
func NewWandfileInstallCommandHandler(
	wandfileRepo interfaces.WandfileRepository,
	wandfileSvc interfaces.WandfileManager,
) *WandfileInstallCommandHandler {
	return &WandfileInstallCommandHandler{
		wandfileRepo: wandfileRepo,
		wandfileSvc:  wandfileSvc,
	}
}

// Handle executes the wandfile install command
func (h *WandfileInstallCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	wandfilePath := "./wandfile"
	if len(args) > 0 {
		wandfilePath = args[0]
	}

	// Check if wandfile exists
	if !h.wandfileRepo.Exists(wandfilePath) {
		return fmt.Errorf("wandfile not found at %s", wandfilePath)
	}

	// Load wandfile
	wandfile, err := h.wandfileRepo.Load(wandfilePath)
	if err != nil {
		return fmt.Errorf("failed to load wandfile: %w", err)
	}

	ctx.Printf("Installing packages from wandfile...\n")

	// Install all packages
	if err := h.wandfileSvc.Install(wandfile); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	ctx.Printf("✓ Successfully installed all packages from wandfile\n")

	return nil
}

// WandfileCheckCommandHandler handles the wandfile check command
type WandfileCheckCommandHandler struct {
	wandfileRepo interfaces.WandfileRepository
	wandfileSvc  interfaces.WandfileManager
}

// NewWandfileCheckCommandHandler creates a new wandfile check command handler
func NewWandfileCheckCommandHandler(
	wandfileRepo interfaces.WandfileRepository,
	wandfileSvc interfaces.WandfileManager,
) *WandfileCheckCommandHandler {
	return &WandfileCheckCommandHandler{
		wandfileRepo: wandfileRepo,
		wandfileSvc:  wandfileSvc,
	}
}

// Handle executes the wandfile check command
func (h *WandfileCheckCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	wandfilePath := "./wandfile"
	if len(args) > 0 {
		wandfilePath = args[0]
	}

	// Check if wandfile exists
	if !h.wandfileRepo.Exists(wandfilePath) {
		return fmt.Errorf("wandfile not found at %s", wandfilePath)
	}

	// Load wandfile
	wandfile, err := h.wandfileRepo.Load(wandfilePath)
	if err != nil {
		return fmt.Errorf("failed to load wandfile: %w", err)
	}

	// Check packages
	missing, err := h.wandfileSvc.Check(wandfile)
	if err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	if len(missing) == 0 {
		ctx.Printf("✓ All packages are installed correctly\n")
		return nil
	}

	ctx.Printf("Missing packages:\n")
	for _, pkg := range missing {
		ctx.Printf("  - %s\n", pkg)
	}
	ctx.Printf("\nRun 'wand wandfile install' to install missing packages\n")

	return nil
}

// WandfileDumpCommandHandler handles the wandfile dump command
type WandfileDumpCommandHandler struct {
	wandfileRepo interfaces.WandfileRepository
	wandfileSvc  interfaces.WandfileManager
}

// NewWandfileDumpCommandHandler creates a new wandfile dump command handler
func NewWandfileDumpCommandHandler(
	wandfileRepo interfaces.WandfileRepository,
	wandfileSvc interfaces.WandfileManager,
) *WandfileDumpCommandHandler {
	return &WandfileDumpCommandHandler{
		wandfileRepo: wandfileRepo,
		wandfileSvc:  wandfileSvc,
	}
}

// Handle executes the wandfile dump command
func (h *WandfileDumpCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	wandfilePath := "./wandfile"
	if len(args) > 0 {
		wandfilePath = args[0]
	}

	// Generate wandfile from installed packages
	wandfile, err := h.wandfileSvc.Dump()
	if err != nil {
		return fmt.Errorf("failed to generate wandfile: %w", err)
	}

	// Save wandfile
	if err := h.wandfileRepo.Save(wandfilePath, wandfile); err != nil {
		return fmt.Errorf("failed to save wandfile: %w", err)
	}

	ctx.Printf("✓ Saved wandfile to %s\n", wandfilePath)

	return nil
}

// SearchCommandHandler handles the search command
type SearchCommandHandler struct {
	formulaRepo interfaces.FormulaRepository
}

// NewSearchCommandHandler creates a new search command handler
func NewSearchCommandHandler(formulaRepo interfaces.FormulaRepository) *SearchCommandHandler {
	return &SearchCommandHandler{
		formulaRepo: formulaRepo,
	}
}

// Handle executes the search command
func (h *SearchCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("search term required")
	}

	searchTerm := strings.ToLower(args[0])

	// Get all formulas
	formulas, err := h.formulaRepo.ListFormulas()
	if err != nil {
		return fmt.Errorf("failed to list formulas: %w", err)
	}

	// Filter matching formulas
	var matches []*entities.Formula
	for _, formula := range formulas {
		nameLower := strings.ToLower(formula.Name)
		descLower := strings.ToLower(formula.Description)

		if strings.Contains(nameLower, searchTerm) || strings.Contains(descLower, searchTerm) {
			matches = append(matches, formula)
		}

		// Also check tags
		for _, tag := range formula.Tags {
			if strings.Contains(strings.ToLower(tag), searchTerm) {
				matches = append(matches, formula)
				break
			}
		}
	}

	if len(matches) == 0 {
		ctx.Printf("No packages found matching '%s'\n", searchTerm)
		return nil
	}

	ctx.Printf("Found %d package(s) matching '%s':\n\n", len(matches), searchTerm)
	for _, formula := range matches {
		ctx.Printf("  %s - %s\n", formula.Name, formula.Description)
		if len(formula.Tags) > 0 {
			ctx.Printf("    Tags: %s\n", strings.Join(formula.Tags, ", "))
		}
	}

	return nil
}

// DoctorCommandHandler handles the doctor command
type DoctorCommandHandler struct {
	registryRepo interfaces.RegistryRepository
	formulaRepo  interfaces.FormulaRepository
	fs           interfaces.FileSystem
	wandDir      string
}

// NewDoctorCommandHandler creates a new doctor command handler
func NewDoctorCommandHandler(
	registryRepo interfaces.RegistryRepository,
	formulaRepo interfaces.FormulaRepository,
	fs interfaces.FileSystem,
	wandDir string,
) *DoctorCommandHandler {
	return &DoctorCommandHandler{
		registryRepo: registryRepo,
		formulaRepo:  formulaRepo,
		fs:           fs,
		wandDir:      wandDir,
	}
}

// Handle executes the doctor command
func (h *DoctorCommandHandler) Handle(ctx interfaces.CommandContext) error {
	ctx.Printf("Running Wand health checks...\n\n")

	allGood := true

	// Check wand directory
	ctx.Printf("Checking Wand installation:\n")
	if h.fs.Exists(h.wandDir) {
		ctx.Printf("  ✓ Wand directory exists: %s\n", h.wandDir)
	} else {
		ctx.Printf("  ✗ Wand directory missing: %s\n", h.wandDir)
		allGood = false
	}

	// Check registry
	if h.registryRepo.Exists() {
		ctx.Printf("  ✓ Package registry exists\n")

		registry, err := h.registryRepo.Load()
		if err != nil {
			ctx.Printf("  ✗ Registry is corrupted: %v\n", err)
			allGood = false
		} else {
			ctx.Printf("  ✓ Registry is valid (%d packages installed)\n", len(registry.Packages))
		}
	} else {
		ctx.Printf("  ⚠ No packages installed yet\n")
	}

	// Check formulas directory
	formulas, err := h.formulaRepo.ListFormulas()
	if err != nil {
		ctx.Printf("  ✗ Cannot access formulas: %v\n", err)
		allGood = false
	} else {
		ctx.Printf("  ✓ Formulas accessible (%d available)\n", len(formulas))
	}

	// Check shims directory
	shimsDir := h.wandDir + "/shims"
	if h.fs.Exists(shimsDir) {
		ctx.Printf("  ✓ Shims directory exists\n")
	} else {
		ctx.Printf("  ⚠ Shims directory not found (will be created on first install)\n")
	}

	ctx.Printf("\n")
	if allGood {
		ctx.Printf("✓ All checks passed!\n")
	} else {
		ctx.Printf("⚠ Some issues detected. Please review the output above.\n")
	}

	return nil
}

// UpdateCommandHandler handles the update command
type UpdateCommandHandler struct {
	installOrchestrator *InstallOrchestrator
	registryRepo        interfaces.RegistryRepository
}

// NewUpdateCommandHandler creates a new update command handler
func NewUpdateCommandHandler(
	installOrchestrator *InstallOrchestrator,
	registryRepo interfaces.RegistryRepository,
) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		installOrchestrator: installOrchestrator,
		registryRepo:        registryRepo,
	}
}

// Handle executes the update command
func (h *UpdateCommandHandler) Handle(ctx interfaces.CommandContext) error {
	args := ctx.GetArgs()
	if len(args) == 0 {
		return fmt.Errorf("package name required")
	}

	packageName := args[0]

	ctx.Printf("Updating %s to latest version...\n", packageName)

	// Install latest version with force flag
	opts := InstallPackageOptions{
		Global: true,
		Force:  true,
	}

	if err := h.installOrchestrator.InstallPackageWithOptions(packageName, "latest", opts); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	ctx.Printf("✓ Successfully updated %s\n", packageName)

	return nil
}

// VersionCommandHandler handles the version command
type VersionCommandHandler struct {
	version   string
	buildTime string
	commit    string
}

// NewVersionCommandHandler creates a new version command handler
func NewVersionCommandHandler(version, buildTime, commit string) *VersionCommandHandler {
	return &VersionCommandHandler{
		version:   version,
		buildTime: buildTime,
		commit:    commit,
	}
}

// Handle executes the version command
func (h *VersionCommandHandler) Handle(ctx interfaces.CommandContext) error {
	ctx.Printf("wand version %s\n", h.version)
	ctx.Printf("Built: %s\n", h.buildTime)
	ctx.Printf("Commit: %s\n", h.commit)
	ctx.Printf("Platform: %s\n", entities.CurrentPlatform().String())
	return nil
}

// OutdatedCommandHandler handles the outdated command
type OutdatedCommandHandler struct {
	installOrchestrator *InstallOrchestrator
	registryRepo        interfaces.RegistryRepository
}

// NewOutdatedCommandHandler creates a new outdated command handler
func NewOutdatedCommandHandler(
	installOrchestrator *InstallOrchestrator,
	registryRepo interfaces.RegistryRepository,
) *OutdatedCommandHandler {
	return &OutdatedCommandHandler{
		installOrchestrator: installOrchestrator,
		registryRepo:        registryRepo,
	}
}

// Handle executes the outdated command
func (h *OutdatedCommandHandler) Handle(ctx interfaces.CommandContext) error {
	registry, err := h.registryRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if len(registry.Packages) == 0 {
		ctx.Printf("No packages installed\n")
		return nil
	}

	ctx.Printf("Checking for outdated packages...\n\n")

	hasOutdated := false
	for name, entry := range registry.Packages {
		// Get current version
		currentVersion, hasGlobal := registry.GetGlobalVersion(name)
		if !hasGlobal && len(entry.Versions) > 0 {
			// Use first version if no global set
			for v := range entry.Versions {
				currentVersion = v
				break
			}
		}

		// Get available versions
		availableVersions, err := h.installOrchestrator.ListAvailableVersions(name)
		if err != nil {
			ctx.Printf("  %s: unable to check for updates\n", name)
			continue
		}

		if len(availableVersions) == 0 {
			continue
		}

		// Get latest version
		latestVersion := availableVersions[0]
		for _, v := range availableVersions {
			if v.Compare(latestVersion) > 0 {
				latestVersion = v
			}
		}

		// Compare versions
		current, err := entities.NewVersion(currentVersion)
		if err != nil {
			continue
		}

		if latestVersion.Compare(current) > 0 {
			ctx.Printf("  %s: %s → %s\n", name, currentVersion, latestVersion.String())
			hasOutdated = true
		}
	}

	if !hasOutdated {
		ctx.Printf("✓ All packages are up to date\n")
	} else {
		ctx.Printf("\nRun 'wand update <package>' to update\n")
	}

	return nil
}
