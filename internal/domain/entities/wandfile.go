package entities

// Wandfile represents a wandfile for declarative system configuration
type Wandfile struct {
	CLI      []WandfileCLI     `yaml:"cli"`      // CLI packages with versions
	GUI      []string          `yaml:"gui"`      // GUI packages (no versions)
	Dotfiles *WandfileDotfiles `yaml:"dotfiles"` // Dotfile configuration
}

// WandfileCLI represents a CLI package entry
type WandfileCLI struct {
	Name    string `yaml:"name"`    // Package name
	Version string `yaml:"version"` // Version constraint or exact version
}

// WandfileDotfiles represents dotfile configuration
type WandfileDotfiles struct {
	Repo     string            `yaml:"repo"`     // Git repository URL
	Symlinks map[string]string `yaml:"symlinks"` // target -> source
}

// NewWandfile creates a new Wandfile
func NewWandfile() *Wandfile {
	return &Wandfile{
		CLI: make([]WandfileCLI, 0),
		GUI: make([]string, 0),
	}
}

// AddCLI adds a CLI package
func (w *Wandfile) AddCLI(name, version string) {
	w.CLI = append(w.CLI, WandfileCLI{
		Name:    name,
		Version: version,
	})
}

// AddGUI adds a GUI package
func (w *Wandfile) AddGUI(name string) {
	w.GUI = append(w.GUI, name)
}

// SetDotfiles sets the dotfiles configuration
func (w *Wandfile) SetDotfiles(repo string, symlinks map[string]string) {
	w.Dotfiles = &WandfileDotfiles{
		Repo:     repo,
		Symlinks: symlinks,
	}
}

// HasCLI checks if a CLI package is defined
func (w *Wandfile) HasCLI(name string) bool {
	for _, cli := range w.CLI {
		if cli.Name == name {
			return true
		}
	}
	return false
}

// HasGUI checks if a GUI package is defined
func (w *Wandfile) HasGUI(name string) bool {
	for _, gui := range w.GUI {
		if gui == name {
			return true
		}
	}
	return false
}

// HasDotfiles checks if dotfiles are configured
func (w *Wandfile) HasDotfiles() bool {
	return w.Dotfiles != nil && w.Dotfiles.Repo != ""
}

// IsEmpty returns true if the wandfile has no entries
func (w *Wandfile) IsEmpty() bool {
	return len(w.CLI) == 0 && len(w.GUI) == 0 && !w.HasDotfiles()
}
