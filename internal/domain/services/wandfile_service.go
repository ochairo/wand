package services

import (
	"fmt"

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
	// Install CLI packages
	for _, cliPkg := range wandfile.CLI {
		if err := s.installerSvc.InstallPackage(cliPkg.Name, cliPkg.Version); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Failed to install %s@%s", cliPkg.Name, cliPkg.Version), err)
		}
	}

	// Install GUI packages
	for _, guiPkg := range wandfile.GUI {
		if err := s.installerSvc.InstallPackage(guiPkg, "latest"); err != nil {
			return errs.Wrap(errs.ErrInstallationFailed, fmt.Sprintf("Failed to install GUI app %s", guiPkg), err)
		}
	}

	// Configure dotfiles if specified
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

	// Check CLI packages
	for _, cliPkg := range wandfile.CLI {
		if _, exists := registry.GetPackage(cliPkg.Name, cliPkg.Version); !exists {
			missing = append(missing, fmt.Sprintf("%s@%s", cliPkg.Name, cliPkg.Version))
		}
	}

	// Check GUI packages
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

	// Add all packages to wandfile
	for name, entry := range registry.Packages {
		// Get global version if set
		globalVersion, hasGlobal := registry.GetGlobalVersion(name)

		switch entry.Type {
		case entities.PackageTypeCLI:
			version := globalVersion
			if !hasGlobal && len(entry.Versions) > 0 {
				// Use any version if no global set
				for v := range entry.Versions {
					version = v
					break
				}
			}
			wandfile.AddCLI(name, version)

		case entities.PackageTypeGUI:
			wandfile.AddGUI(name)
		}
	}

	// Add dotfiles config if exists
	if s.dotfileRepo.Exists() {
		dotfileConfig, err := s.dotfileRepo.Load()
		if err == nil && dotfileConfig.HasSymlinks() {
			wandfile.SetDotfiles(dotfileConfig.RepoURL, dotfileConfig.Symlinks)
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
