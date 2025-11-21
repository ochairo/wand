package types

import "time"

// Registry represents the local package registry
type Registry struct {
	Packages       map[string]*PackageEntry // name -> PackageEntry
	GlobalVersions map[string]string        // name -> version
	UpdatedAt      time.Time
}

// PackageEntry represents all installed versions of a package
type PackageEntry struct {
	Name     string
	Type     PackageType
	Versions map[string]*Package // version -> Package
}
