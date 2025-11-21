package entities

import "time"

// Registry represents the local package registry
type Registry struct {
	Packages       map[string]*PackageEntry `json:"packages"`        // name -> PackageEntry
	GlobalVersions map[string]string        `json:"global_versions"` // name -> version
	UpdatedAt      time.Time                `json:"updated_at"`
}

// PackageEntry represents all installed versions of a package
type PackageEntry struct {
	Name     string              `json:"name"`
	Type     PackageType         `json:"type"`
	Versions map[string]*Package `json:"versions"` // version -> Package
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		Packages:       make(map[string]*PackageEntry),
		GlobalVersions: make(map[string]string),
		UpdatedAt:      time.Now(),
	}
}

// AddPackage adds a package to the registry
func (r *Registry) AddPackage(pkg *Package) {
	entry, exists := r.Packages[pkg.Name]
	if !exists {
		entry = &PackageEntry{
			Name:     pkg.Name,
			Type:     pkg.Type,
			Versions: make(map[string]*Package),
		}
		r.Packages[pkg.Name] = entry
	}

	entry.Versions[pkg.VersionString()] = pkg
	r.UpdatedAt = time.Now()
}

// RemovePackage removes a specific version of a package
func (r *Registry) RemovePackage(name, version string) bool {
	entry, exists := r.Packages[name]
	if !exists {
		return false
	}

	if _, ok := entry.Versions[version]; !ok {
		return false
	}

	delete(entry.Versions, version)

	// Remove entry if no versions left
	if len(entry.Versions) == 0 {
		delete(r.Packages, name)
	}

	// Clear global version if it was removed
	if r.GlobalVersions[name] == version {
		delete(r.GlobalVersions, name)
	}

	r.UpdatedAt = time.Now()
	return true
}

// GetPackage returns a specific package version
func (r *Registry) GetPackage(name, version string) (*Package, bool) {
	entry, exists := r.Packages[name]
	if !exists {
		return nil, false
	}

	pkg, ok := entry.Versions[version]
	return pkg, ok
}

// GetAllVersions returns all installed versions of a package
func (r *Registry) GetAllVersions(name string) ([]*Package, bool) {
	entry, exists := r.Packages[name]
	if !exists {
		return nil, false
	}

	versions := make([]*Package, 0, len(entry.Versions))
	for _, pkg := range entry.Versions {
		versions = append(versions, pkg)
	}

	return versions, true
}

// SetGlobalVersion sets the global version for a package
func (r *Registry) SetGlobalVersion(name, version string) {
	r.GlobalVersions[name] = version
	r.UpdatedAt = time.Now()
}

// GetGlobalVersion returns the global version for a package
func (r *Registry) GetGlobalVersion(name string) (string, bool) {
	version, ok := r.GlobalVersions[name]
	return version, ok
}

// HasPackage checks if a package exists
func (r *Registry) HasPackage(name string) bool {
	_, exists := r.Packages[name]
	return exists
}

// ListAllPackages returns all package entries
func (r *Registry) ListAllPackages() []*PackageEntry {
	entries := make([]*PackageEntry, 0, len(r.Packages))
	for _, entry := range r.Packages {
		entries = append(entries, entry)
	}
	return entries
}
