package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// PackageNameValidator validates package names according to naming conventions.
type PackageNameValidator struct {
	namePattern *regexp.Regexp
}

// NewPackageNameValidator creates a new package name validator.
func NewPackageNameValidator() *PackageNameValidator {
	return &PackageNameValidator{
		namePattern: regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]$"),
	}
}

// Validate checks if a package name meets all validation requirements.
func (v *PackageNameValidator) Validate(name string) error {
	if name == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if !v.namePattern.MatchString(name) {
		return fmt.Errorf("invalid package name")
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("no consecutive hyphens")
	}
	return nil
}
