// Package interfaces defines all domain interfaces following the dependency inversion principle.
// The package name "interfaces" is a standard Go pattern for abstract contracts in domain-driven design.
package interfaces

import "github.com/ochairo/wand/internal/domain/entities"

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string
	Assets  []GitHubAsset
}

// GitHubAsset represents a release asset
type GitHubAsset struct {
	Name        string
	DownloadURL string
}

// GitHubClient defines the interface for GitHub API operations
type GitHubClient interface {
	GetLatestRelease(owner, repo string) (*GitHubRelease, error)
	GetRelease(owner, repo, tag string) (*GitHubRelease, error)
	ListReleases(owner, repo string) ([]*GitHubRelease, error)
	DownloadAsset(asset *GitHubAsset, destPath string) error
}

// PackageResolver defines the interface for resolving package versions
type PackageResolver interface {
	ResolveVersion(pkg *entities.Formula, constraint string) (*entities.Version, error)
	GetAvailableVersions(pkg *entities.Formula) ([]*entities.Version, error)
	GetLatestVersion(pkg *entities.Formula) (*entities.Version, error)
}

// VersionFinder defines the interface for finding the active version
type VersionFinder interface {
	FindVersion(packageName, currentDir string) (string, error)
	GetGlobalVersion(packageName string) (string, error)
}
