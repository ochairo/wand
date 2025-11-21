// Package entities defines the core domain entities.
package entities

// DotfileConfig represents the dotfile repository configuration
type DotfileConfig struct {
	RepoURL  string            `json:"repo_url"`  // Git repository URL
	LocalDir string            `json:"local_dir"` // Local checkout directory
	Symlinks map[string]string `json:"symlinks"`  // target -> source mapping
}

// NewDotfileConfig creates a new DotfileConfig
func NewDotfileConfig(repoURL, localDir string) *DotfileConfig {
	return &DotfileConfig{
		RepoURL:  repoURL,
		LocalDir: localDir,
		Symlinks: make(map[string]string),
	}
}

// AddSymlink adds a symlink mapping
func (d *DotfileConfig) AddSymlink(target, source string) {
	d.Symlinks[target] = source
}

// RemoveSymlink removes a symlink mapping
func (d *DotfileConfig) RemoveSymlink(target string) bool {
	if _, ok := d.Symlinks[target]; ok {
		delete(d.Symlinks, target)
		return true
	}
	return false
}

// GetSource returns the source for a given target
func (d *DotfileConfig) GetSource(target string) (string, bool) {
	source, ok := d.Symlinks[target]
	return source, ok
}

// HasSymlinks returns true if there are any symlink mappings
func (d *DotfileConfig) HasSymlinks() bool {
	return len(d.Symlinks) > 0
}
