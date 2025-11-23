package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	errs "github.com/ochairo/wand/internal/domain/errors"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// WandfileService handles wandfile operations
type WandfileService struct {
	wandfileRepo  interfaces.WandfileRepository
	registryRepo  interfaces.RegistryRepository
	installerSvc  *InstallerService
	versionSvc    *VersionService
	dotfileRepo   interfaces.DotfileRepository
	fs            interfaces.FileSystem
	shellExecutor interfaces.ShellExecutor
	homeDir       string
}

// NewWandfileService creates a new wandfile service
func NewWandfileService(
	wandfileRepo interfaces.WandfileRepository,
	registryRepo interfaces.RegistryRepository,
	installerSvc *InstallerService,
	versionSvc *VersionService,
	dotfileRepo interfaces.DotfileRepository,
	fs interfaces.FileSystem,
	shellExecutor interfaces.ShellExecutor,
	homeDir string,
) *WandfileService {
	return &WandfileService{
		wandfileRepo:  wandfileRepo,
		registryRepo:  registryRepo,
		installerSvc:  installerSvc,
		versionSvc:    versionSvc,
		dotfileRepo:   dotfileRepo,
		fs:            fs,
		shellExecutor: shellExecutor,
		homeDir:       homeDir,
	}
}

// Install installs all packages and configures dotfiles from a wandfile
func (s *WandfileService) Install(wandfile *entities.Wandfile) error {
	// Install formulas (new format)
	for _, formula := range wandfile.Formulas {
		name, version := s.parseFormula(formula)
		if err := s.installerSvc.InstallPackage(name, version); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Failed to install %s@%s", name, version), err)
		}
	}

	// Install CLI packages (legacy format)
	for _, cliPkg := range wandfile.CLI {
		if err := s.installerSvc.InstallPackage(cliPkg.Name, cliPkg.Version); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Failed to install %s@%s", cliPkg.Name, cliPkg.Version), err)
		}
	}

	// Install GUI packages (legacy format)
	for _, guiPkg := range wandfile.GUI {
		if err := s.installerSvc.InstallPackage(guiPkg, "latest"); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Failed to install GUI app %s", guiPkg), err)
		}
	}

	// Configure symlinks (new format)
	if wandfile.HasSymlinks() {
		if err := s.configureSymlinks(wandfile.Symlinks); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, "Failed to configure symlinks", err)
		}
	}

	// Configure dotfiles (legacy format)
	if wandfile.HasDotfiles() {
		if err := s.configureDotfiles(wandfile.Dotfiles); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, "Failed to configure dotfiles", err)
		}
	}

	return nil
}

// Check verifies that all packages in wandfile are installed correctly
func (s *WandfileService) Check(wandfile *entities.Wandfile) ([]string, error) {
	var missing []string

	registry, err := s.registryRepo.Load()
	if err != nil {
		return nil, errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry", err)
	}

	// Check formulas (new format)
	for _, formula := range wandfile.Formulas {
		name, version := s.parseFormula(formula)
		if _, exists := registry.GetPackage(name, version); !exists {
			if !registry.HasPackage(name) {
				missing = append(missing, fmt.Sprintf("%s@%s", name, version))
			}
		}
	}

	// Check CLI packages (legacy)
	for _, cliPkg := range wandfile.CLI {
		if _, exists := registry.GetPackage(cliPkg.Name, cliPkg.Version); !exists {
			missing = append(missing, fmt.Sprintf("%s@%s", cliPkg.Name, cliPkg.Version))
		}
	}

	// Check GUI packages (legacy)
	for _, guiPkg := range wandfile.GUI {
		if !registry.HasPackage(guiPkg) {
			missing = append(missing, guiPkg)
		}
	}

	return missing, nil
}

// Dump generates a wandfile from currently installed packages
func (s *WandfileService) Dump() (*entities.Wandfile, error) {
	registry, err := s.registryRepo.Load()
	if err != nil {
		return nil, errs.Wrap(errs.ErrRegistryCorrupted, "Failed to load registry", err)
	}

	wandfile := entities.NewWandfile()

	// Add all packages to wandfile using new formulas format
	for name, entry := range registry.Packages {
		// Get global version if set
		globalVersion, hasGlobal := registry.GetGlobalVersion(name)

		version := globalVersion
		if !hasGlobal && len(entry.Versions) > 0 {
			// Use any version if no global set
			for v := range entry.Versions {
				version = v
				break
			}
		}

		// Add to formulas
		wandfile.AddFormula(name, version)
	}

	// Add symlinks if dotfile config exists
	if s.dotfileRepo.Exists() {
		dotfileConfig, err := s.dotfileRepo.Load()
		if err == nil && dotfileConfig.HasSymlinks() {
			for target, source := range dotfileConfig.Symlinks {
				wandfile.AddSymlink(target, source)
			}
		}
	}

	return wandfile, nil
}

// Update updates all packages in wandfile to their latest versions
func (s *WandfileService) Update() error {
	// Load wandfile from home directory
	wandfile, err := s.wandfileRepo.Load(s.homeDir)
	if err != nil {
		return errs.Wrap(errs.ErrFileNotFound, "Failed to load wandfile", err)
	}

	if wandfile == nil {
		return errs.New(errs.ErrFileNotFound, "Wandfile not found in home directory")
	}

	// Track updates
	updated := 0
	skipped := 0

	// Update CLI packages
	for i, cliPkg := range wandfile.CLI {
		// Get latest version
		latestVersion, err := s.versionSvc.ResolveVersion(cliPkg.Name, "latest")
		if err != nil {
			fmt.Printf("⚠ Skipped %s: %v\n", cliPkg.Name, err)
			skipped++
			continue
		}

		if latestVersion.String() != cliPkg.Version {
			fmt.Printf("✓ %s: %s → %s\n", cliPkg.Name, cliPkg.Version, latestVersion.String())
			wandfile.CLI[i].Version = latestVersion.String()
			updated++

			// Install updated version
			if err := s.installerSvc.InstallPackage(cliPkg.Name, latestVersion.String()); err != nil {
				fmt.Printf("⚠ Failed to install %s@%s: %v\n", cliPkg.Name, latestVersion.String(), err)
			}
		}
	}

	// Update GUI packages
	for _, guiName := range wandfile.GUI {
		// Get latest version
		latestVersion, err := s.versionSvc.ResolveVersion(guiName, "latest")
		if err != nil {
			fmt.Printf("⚠ Skipped %s: %v\n", guiName, err)
			skipped++
			continue
		}

		fmt.Printf("✓ %s: installed (GUI packages check only)\n", guiName)
		updated++

		// Install updated version
		if err := s.installerSvc.InstallPackage(guiName, latestVersion.String()); err != nil {
			fmt.Printf("⚠ Failed to install %s@%s: %v\n", guiName, latestVersion.String(), err)
		}
	}

	// Save updated wandfile
	if err := s.wandfileRepo.Save(s.homeDir, wandfile); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to save wandfile", err)
	}

	fmt.Printf("\n✨ Update complete: %d updated, %d skipped\n", updated, skipped)
	return nil
}

// configureDotfiles sets up dotfile repository and symlinks
func (s *WandfileService) configureDotfiles(dotfiles *entities.WandfileDotfiles) error {
	// Create dotfile config
	config := entities.NewDotfileConfig(dotfiles.Repo, s.homeDir+"/.dotfiles")

	// Add all symlinks
	for target, source := range dotfiles.Symlinks {
		config.AddSymlink(target, source)
	}

	// Save dotfile config
	if err := s.dotfileRepo.Save(config); err != nil {
		return errs.Wrap(errs.ErrConfigInvalid, "Failed to save dotfile config", err)
	}

	// Clone dotfile repository if it doesn't exist
	if !s.fs.Exists(config.LocalDir) {
		if _, err := s.shellExecutor.Execute("git", "clone", config.RepoURL, config.LocalDir); err != nil {
			return errs.Wrap(errs.ErrNetworkUnreachable, "Failed to clone dotfile repository", err)
		}
	}

	// Create symlinks
	for target, source := range config.Symlinks {
		targetPath := s.homeDir + "/" + target
		sourcePath := config.LocalDir + "/" + source

		// Check if source exists
		if !s.fs.Exists(sourcePath) {
			return errs.NewWithDetails(errs.ErrFileNotFound, "Dotfile source not found", fmt.Sprintf("path: %q", sourcePath))
		}

		// Backup existing file if it exists and is not a symlink
		if s.fs.Exists(targetPath) {
			link, err := s.fs.ReadSymlink(targetPath)
			if err != nil {
				// Not a symlink, backup the file
				backupPath := targetPath + ".wand-backup"
				if _, err := s.shellExecutor.Execute("mv", targetPath, backupPath); err != nil {
					return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to backup %s", targetPath), err)
				}
			} else {
				// Remove existing symlink
				if link != sourcePath {
					if err := s.fs.Remove(targetPath); err != nil {
						return errs.Wrap(errs.ErrPermissionDenied, "Failed to remove old symlink", err)
					}
				} else {
					// Already correctly linked
					continue
				}
			}
		}

		// Create symlink
		if err := s.fs.Symlink(sourcePath, targetPath); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to create symlink %s -> %s", targetPath, sourcePath), err)
		}
	}

	return nil
}

// configureSymlinks sets up symlinks without requiring a dotfile repo
func (s *WandfileService) configureSymlinks(symlinks map[string]string) error {
	// Assume dotfiles are in ~/.dotfiles by default
	dotfilesDir := filepath.Join(s.homeDir, ".dotfiles")

	// Create symlinks directly from wandfile
	for target, source := range symlinks {
		targetPath := filepath.Join(s.homeDir, target)
		sourcePath := filepath.Join(dotfilesDir, source)

		// Check if source exists
		if !s.fs.Exists(sourcePath) {
			return errs.NewWithDetails(errs.ErrFileNotFound, "Symlink source not found", fmt.Sprintf("path: %q (expected in ~/.dotfiles/)", sourcePath))
		}

		// Ensure target directory exists
		targetDir := filepath.Dir(targetPath)
		if !s.fs.Exists(targetDir) {
			if _, err := s.shellExecutor.Execute("mkdir", "-p", targetDir); err != nil {
				return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to create directory %s", targetDir), err)
			}
		}

		// Backup existing file if it exists and is not a symlink
		if s.fs.Exists(targetPath) {
			link, err := s.fs.ReadSymlink(targetPath)
			if err != nil {
				// Not a symlink, backup the file
				backupPath := targetPath + ".wand-backup"
				if _, err := s.shellExecutor.Execute("mv", targetPath, backupPath); err != nil {
					return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to backup %s", targetPath), err)
				}
			} else {
				// Remove existing symlink
				if link != sourcePath {
					if err := s.fs.Remove(targetPath); err != nil {
						return errs.Wrap(errs.ErrPermissionDenied, "Failed to remove old symlink", err)
					}
				} else {
					// Already correctly linked
					continue
				}
			}
		}

		// Create symlink
		if err := s.fs.Symlink(sourcePath, targetPath); err != nil {
			return errs.Wrap(errs.ErrPermissionDenied, fmt.Sprintf("Failed to create symlink %s -> %s", targetPath, sourcePath), err)
		}
	}

	return nil
}

// ValidationResult represents wandfile validation result
type ValidationResult struct {
	Valid        bool
	Errors       []string
	Packages     []string
	HasDotfiles  bool
	SymlinkCount int
}

// WandfileSummary represents wandfile summary for display
type WandfileSummary struct {
	CLIPackages  []string
	GUIApps      []string
	HasDotfiles  bool
	DotfilesRepo string
	SymlinkCount int
}

// InstallFrom installs from a wandfile (path, URL, or GitHub shorthand)
func (s *WandfileService) InstallFrom(source string) error {
	// TODO: Handle remote URLs and GitHub shorthands
	// For now, just handle local paths

	wandfile, err := s.wandfileRepo.Load(source)
	if err != nil {
		return errs.Wrap(errs.ErrFileNotFound, "Failed to load wandfile", err)
	}

	return s.Install(wandfile)
}

// Init creates a new wandfile interactively
func (s *WandfileService) Init(ctx interfaces.CommandContext) error {
	// Create basic wandfile template
	wandfile := entities.NewWandfile()

	// Add some common packages
	wandfile.AddCLI("git", "latest")
	wandfile.AddCLI("zsh", "latest")
	wandfile.AddCLI("starship", "latest")
	wandfile.AddCLI("neovim", "latest")
	wandfile.AddCLI("bat", "latest")
	wandfile.AddCLI("eza", "latest")
	wandfile.AddCLI("ripgrep", "latest")
	wandfile.AddCLI("fzf", "latest")

	// Save to current directory
	if err := s.wandfileRepo.Save(".", wandfile); err != nil {
		return errs.Wrap(errs.ErrPermissionDenied, "Failed to save wandfile", err)
	}

	return nil
}

// Validate validates a wandfile
func (s *WandfileService) Validate(path string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]string, 0),
		Packages: make([]string, 0),
	}

	wandfile, err := s.wandfileRepo.Load(path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to load wandfile: %v", err))
		return result, nil
	}

	// Check formulas (new format)
	for _, formula := range wandfile.Formulas {
		name, version := s.parseFormula(formula)
		result.Packages = append(result.Packages, fmt.Sprintf("%s@%s", name, version))
		// TODO: Check if formula exists
	}

	// Check CLI packages (legacy)
	for _, cli := range wandfile.CLI {
		result.Packages = append(result.Packages, fmt.Sprintf("%s@%s", cli.Name, cli.Version))
		// TODO: Check if formula exists
	}

	// Check GUI apps (legacy)
	result.Packages = append(result.Packages, wandfile.GUI...)

	// Check symlinks (new format)
	if wandfile.HasSymlinks() {
		result.SymlinkCount = len(wandfile.Symlinks)
	}

	// Check dotfiles (legacy)
	if wandfile.HasDotfiles() {
		result.HasDotfiles = true
		result.SymlinkCount += len(wandfile.Dotfiles.Symlinks)
		// TODO: Validate repo URL is accessible
	}

	return result, nil
}

// DumpToString generates wandfile content from current system
func (s *WandfileService) DumpToString() (string, error) {
	wandfile, err := s.Dump()
	if err != nil {
		return "", err
	}

	// Convert to YAML string
	content := "# Complete development environment configuration\n"
	content += "# Install with: wand wandfile install\n\n"

	// Formulas section (new format)
	if len(wandfile.Formulas) > 0 {
		content += "formulas:\n"
		for _, formula := range wandfile.Formulas {
			content += fmt.Sprintf("  - %s\n", formula)
		}
	}

	// Symlinks section (new format)
	if wandfile.HasSymlinks() {
		content += "\nsymlinks:\n"
		for target, source := range wandfile.Symlinks {
			content += fmt.Sprintf("  %s: %s\n", target, source)
		}
	}

	return content, nil
}

// Show displays a summary of a wandfile
func (s *WandfileService) Show(path string) (*WandfileSummary, error) {
	wandfile, err := s.wandfileRepo.Load(path)
	if err != nil {
		return nil, errs.Wrap(errs.ErrFileNotFound, "Failed to load wandfile", err)
	}

	summary := &WandfileSummary{
		CLIPackages: make([]string, 0),
		GUIApps:     make([]string, 0),
	}

	// Add formulas (new format)
	for _, formula := range wandfile.Formulas {
		name, version := s.parseFormula(formula)
		summary.CLIPackages = append(summary.CLIPackages, fmt.Sprintf("%s@%s", name, version))
	}

	// Add CLI packages (legacy)
	for _, cli := range wandfile.CLI {
		summary.CLIPackages = append(summary.CLIPackages, fmt.Sprintf("%s@%s", cli.Name, cli.Version))
	}

	// Add GUI apps (legacy)
	summary.GUIApps = append(summary.GUIApps, wandfile.GUI...)

	// Check symlinks (new format)
	if wandfile.HasSymlinks() {
		summary.SymlinkCount = len(wandfile.Symlinks)
	}

	// Check dotfiles (legacy)
	if wandfile.HasDotfiles() {
		summary.HasDotfiles = true
		summary.DotfilesRepo = wandfile.Dotfiles.Repo
		summary.SymlinkCount += len(wandfile.Dotfiles.Symlinks)
	}

	return summary, nil
}

// CheckPath checks wandfile against installed packages
func (s *WandfileService) CheckPath(path string) ([]string, error) {
	wandfile, err := s.wandfileRepo.Load(path)
	if err != nil {
		return nil, errs.Wrap(errs.ErrFileNotFound, "Failed to load wandfile", err)
	}

	return s.Check(wandfile)
}

// parseFormula parses a formula string in format "name" or "name@version"
func (s *WandfileService) parseFormula(formula string) (name, version string) {
	parts := strings.Split(formula, "@")
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	} else {
		version = "latest"
	}
	return name, version
}
