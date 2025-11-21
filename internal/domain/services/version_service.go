package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ochairo/wand/internal/domain/entities"
	errs "github.com/ochairo/wand/internal/domain/errors"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// VersionService handles version resolution and comparison logic
type VersionService struct {
	githubClient interfaces.GitHubClient
	formulaRepo  interfaces.FormulaRepository
}

// NewVersionService creates a new version service
func NewVersionService(
	githubClient interfaces.GitHubClient,
	formulaRepo interfaces.FormulaRepository,
) *VersionService {
	return &VersionService{
		githubClient: githubClient,
		formulaRepo:  formulaRepo,
	}
}

// parseRepository splits "owner/repo" into owner and repo
func parseRepository(repository string) (string, string, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return "", "", errs.NewWithDetails(errs.ErrInvalidPath, "Invalid repository format", fmt.Sprintf("expected owner/repo, got: %q", repository))
	}
	return parts[0], parts[1], nil
}

// ResolveVersion resolves "latest" or validates specific version
func (s *VersionService) ResolveVersion(packageName, versionStr string) (*entities.Version, error) {
	if versionStr == "" || versionStr == "latest" {
		return s.GetLatestVersion(packageName)
	}

	// Parse and validate the version
	version, err := entities.NewVersion(versionStr)
	if err != nil {
		return nil, errs.New(errs.ErrInvalidVersion, fmt.Sprintf("Invalid version: %q", versionStr))
	}

	// Verify version exists for the package
	exists, err := s.VersionExists(packageName, version)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.NewWithDetails(errs.ErrVersionNotFound, "Version not available", fmt.Sprintf("package: %q, version: %q", packageName, versionStr))
	}

	return version, nil
}

// GetLatestVersion fetches the latest available version for a package
func (s *VersionService) GetLatestVersion(packageName string) (*entities.Version, error) {
	formula, err := s.formulaRepo.GetFormula(packageName)
	if err != nil {
		return nil, errs.NewWithDetails(errs.ErrPackageNotFound, "Formula not found", fmt.Sprintf("package: %q", packageName))
	}

	// Get all releases from GitHub
	owner, repo, err := parseRepository(formula.Repository)
	if err != nil {
		return nil, err
	}

	releases, err := s.githubClient.ListReleases(owner, repo)
	if err != nil {
		return nil, errs.Wrap(errs.ErrNetworkUnreachable, fmt.Sprintf("Failed to fetch releases for %s", packageName), err)
	}

	if len(releases) == 0 {
		return nil, errs.NewWithDetails(errs.ErrVersionNotFound, "No versions found", fmt.Sprintf("package: %q", packageName))
	}

	// Parse versions from releases
	versions := make([]*entities.Version, 0)
	for _, release := range releases {
		// Strip package name prefix from tag (e.g., "nano-8.7" -> "8.7")
		tagName := strings.TrimPrefix(release.TagName, packageName+"-")

		version, err := entities.NewVersion(tagName)
		if err != nil {
			// Skip invalid version tags
			continue
		}
		versions = append(versions, version)
	}

	if len(versions) == 0 {
		return nil, errs.NewWithDetails(errs.ErrVersionNotFound, "No valid versions found", fmt.Sprintf("package: %q", packageName))
	}

	// Sort and return latest
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})

	return versions[0], nil
}

// VersionExists checks if a specific version exists for a package
func (s *VersionService) VersionExists(packageName string, version *entities.Version) (bool, error) {
	formula, err := s.formulaRepo.GetFormula(packageName)
	if err != nil {
		return false, errs.NewWithDetails(errs.ErrPackageNotFound, "Formula not found", fmt.Sprintf("package: %q", packageName))
	}

	owner, repo, err := parseRepository(formula.Repository)
	if err != nil {
		return false, err
	}

	releases, err := s.githubClient.ListReleases(owner, repo)
	if err != nil {
		return false, errs.Wrap(errs.ErrNetworkUnreachable, fmt.Sprintf("Failed to fetch releases for %s", packageName), err)
	}

	// Check if version exists in releases
	versionStr := version.String()
	for _, release := range releases {
		// Strip package name prefix from tag (e.g., "nano-8.7" -> "8.7")
		tagName := strings.TrimPrefix(release.TagName, packageName+"-")

		parsedTag, err := entities.NewVersion(tagName)
		if err != nil {
			continue
		}
		if parsedTag.String() == versionStr {
			return true, nil
		}
	}

	return false, nil
}

// ListAvailableVersions returns all available versions for a package
func (s *VersionService) ListAvailableVersions(packageName string) ([]*entities.Version, error) {
	formula, err := s.formulaRepo.GetFormula(packageName)
	if err != nil {
		return nil, errs.NewWithDetails(errs.ErrPackageNotFound, "Formula not found", fmt.Sprintf("package: %q", packageName))
	}

	owner, repo, err := parseRepository(formula.Repository)
	if err != nil {
		return nil, err
	}

	releases, err := s.githubClient.ListReleases(owner, repo)
	if err != nil {
		return nil, errs.Wrap(errs.ErrNetworkUnreachable, fmt.Sprintf("Failed to fetch releases for %s", packageName), err)
	}

	versions := make([]*entities.Version, 0)
	for _, release := range releases {
		// Strip package name prefix from tag (e.g., "nano-8.7" -> "8.7")
		tagName := strings.TrimPrefix(release.TagName, packageName+"-")

		version, err := entities.NewVersion(tagName)
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}

	// Sort versions descending (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})

	return versions, nil
}

// FindBestMatch finds the best version matching a constraint
func (s *VersionService) FindBestMatch(packageName, constraint string) (*entities.Version, error) {
	versions, err := s.ListAvailableVersions(packageName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, errs.NewWithDetails(errs.ErrVersionNotFound, "No versions found", fmt.Sprintf("package: %q", packageName))
	}

	// Handle constraint matching
	switch constraint {
	case "", "latest", "*":
		return versions[0], nil // Already sorted, first is latest
	default:
		// Try semver constraint matching
		if matched := s.matchConstraint(versions, constraint); matched != nil {
			return matched, nil
		}

		// Fallback: try exact match
		targetVersion, err := entities.NewVersion(constraint)
		if err != nil {
			return nil, errs.New(errs.ErrInvalidVersion, fmt.Sprintf("Invalid version constraint: %q", constraint))
		}

		for _, v := range versions {
			if v.Equal(targetVersion) {
				return v, nil
			}
		}

		return nil, errs.NewWithDetails(errs.ErrVersionNotFound, "No version matching constraint", fmt.Sprintf("package: %q, constraint: %q", packageName, constraint))
	}
}

// matchConstraint matches versions against a constraint pattern
// Supports: ^1.2.3 (compatible), ~1.2.3 (approximately), >=1.0.0, <=2.0.0, 1.x, etc.
func (s *VersionService) matchConstraint(versions []*entities.Version, constraint string) *entities.Version {
	if len(constraint) == 0 {
		return nil
	}

	// Handle caret (^) - compatible version
	if constraint[0] == '^' {
		return s.matchCaret(versions, constraint[1:])
	}

	// Handle tilde (~) - approximately version
	if constraint[0] == '~' {
		return s.matchTilde(versions, constraint[1:])
	}

	// Handle comparison operators
	if len(constraint) > 1 {
		switch constraint[:2] {
		case ">=":
			return s.matchGreaterOrEqual(versions, constraint[2:])
		case "<=":
			return s.matchLessOrEqual(versions, constraint[2:])
		}
	}

	if len(constraint) > 0 {
		if constraint[0] == '>' {
			return s.matchGreater(versions, constraint[1:])
		}
		if constraint[0] == '<' {
			return s.matchLess(versions, constraint[1:])
		}
	}

	// Handle x-ranges (1.x, 1.2.x, etc.)
	if strings.Contains(constraint, "x") {
		return s.matchXRange(versions, constraint)
	}

	return nil
}

// matchCaret handles ^ (caret) constraint: ^1.2.3 allows >=1.2.3 <2.0.0
func (s *VersionService) matchCaret(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	for _, v := range versions {
		// Same major version, greater or equal to base
		if v.Major == baseVersion.Major && (v.Minor > baseVersion.Minor ||
			(v.Minor == baseVersion.Minor && v.Patch >= baseVersion.Patch)) {
			return v
		}
	}
	return nil
}

// matchTilde handles ~ (tilde) constraint: ~1.2.3 allows >=1.2.3 <1.3.0
func (s *VersionService) matchTilde(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	for _, v := range versions {
		// Same major.minor, patch >= base patch
		if v.Major == baseVersion.Major && v.Minor == baseVersion.Minor && v.Patch >= baseVersion.Patch {
			return v
		}
	}
	return nil
}

// matchGreaterOrEqual handles >= constraint
func (s *VersionService) matchGreaterOrEqual(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	for _, v := range versions {
		if !v.LessThan(baseVersion) {
			return v
		}
	}
	return nil
}

// matchLessOrEqual handles <= constraint
func (s *VersionService) matchLessOrEqual(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	// Search backwards since versions are sorted descending
	for i := len(versions) - 1; i >= 0; i-- {
		if !versions[i].GreaterThan(baseVersion) {
			return versions[i]
		}
	}
	return nil
}

// matchGreater handles > constraint
func (s *VersionService) matchGreater(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	for _, v := range versions {
		if v.GreaterThan(baseVersion) {
			return v
		}
	}
	return nil
}

// matchLess handles < constraint
func (s *VersionService) matchLess(versions []*entities.Version, constraint string) *entities.Version {
	baseVersion, err := entities.NewVersion(constraint)
	if err != nil {
		return nil
	}

	// Search backwards since versions are sorted descending
	for i := len(versions) - 1; i >= 0; i-- {
		if versions[i].LessThan(baseVersion) {
			return versions[i]
		}
	}
	return nil
}

// matchXRange handles x-ranges: 1.x, 1.2.x, etc.
func (s *VersionService) matchXRange(versions []*entities.Version, constraint string) *entities.Version {
	parts := strings.Split(constraint, ".")
	if len(parts) == 0 {
		return nil
	}

	var major, minor int
	if major, _ = strconv.Atoi(parts[0]); len(parts) == 0 || parts[0] == "x" {
		return versions[0] // Any version
	}

	if len(parts) > 1 && parts[1] != "x" {
		minor, _ = strconv.Atoi(parts[1])
	} else {
		// Match major.x - return latest with matching major
		for _, v := range versions {
			if v.Major == major {
				return v
			}
		}
		return nil
	}

	// Match major.minor.x - return latest with matching major.minor
	for _, v := range versions {
		if v.Major == major && v.Minor == minor {
			return v
		}
	}
	return nil
}

// CompareVersions compares two version strings
func (s *VersionService) CompareVersions(v1Str, v2Str string) (int, error) {
	v1, err := entities.NewVersion(v1Str)
	if err != nil {
		return 0, errs.New(errs.ErrInvalidVersion, fmt.Sprintf("Invalid version: %q", v1Str))
	}

	v2, err := entities.NewVersion(v2Str)
	if err != nil {
		return 0, errs.New(errs.ErrInvalidVersion, fmt.Sprintf("Invalid version: %q", v2Str))
	}

	if v1.GreaterThan(v2) {
		return 1, nil
	} else if v1.LessThan(v2) {
		return -1, nil
	}
	return 0, nil
}
