package validation

import (
	"fmt"
	"regexp"
)

// VersionValidator validates semantic version strings.
type VersionValidator struct {
	versionPattern *regexp.Regexp
}

// NewVersionValidator creates a new version validator.
func NewVersionValidator() *VersionValidator {
	return &VersionValidator{
		versionPattern: regexp.MustCompile(`^[0-9]+\.[0-9]+(\.[0-9]+)?(-[a-z0-9]+)?$`),
	}
}

// Validate checks if a version string is valid.
func (v *VersionValidator) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	if len(version) > 50 {
		return fmt.Errorf("version too long")
	}
	if !v.versionPattern.MatchString(version) {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
