package entities

// WandRC represents a .wandrc file for per-project version overrides
type WandRC struct {
	Versions map[string]string `yaml:"versions"` // package name -> version
}

// NewWandRC creates a new WandRC
func NewWandRC() *WandRC {
	return &WandRC{
		Versions: make(map[string]string),
	}
}

// SetVersion sets a version override for a package
func (w *WandRC) SetVersion(name, version string) {
	w.Versions[name] = version
}

// GetVersion gets the version override for a package
func (w *WandRC) GetVersion(name string) (string, bool) {
	version, ok := w.Versions[name]
	return version, ok
}

// RemoveVersion removes a version override
func (w *WandRC) RemoveVersion(name string) bool {
	if _, ok := w.Versions[name]; ok {
		delete(w.Versions, name)
		return true
	}
	return false
}

// HasVersion checks if a version override exists
func (w *WandRC) HasVersion(name string) bool {
	_, ok := w.Versions[name]
	return ok
}

// IsEmpty returns true if no version overrides are defined
func (w *WandRC) IsEmpty() bool {
	return len(w.Versions) == 0
}
