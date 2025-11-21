// Package validation provides input validation utilities.
package validation

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FormulaValidator validates formula YAML files
type FormulaValidator struct {
	nameValidator    *PackageNameValidator
	versionValidator *VersionValidator
	urlValidator     *URLValidator
	releaseValidator *URLValidator
}

// FormulaSchema represents the expected structure of a formula file
type FormulaSchema struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Homepage    string            `yaml:"homepage"`
	License     string            `yaml:"license"`
	Tags        []string          `yaml:"tags"`
	Releases    map[string]string `yaml:"releases"`
}

// NewFormulaValidator creates a new formula validator
func NewFormulaValidator() *FormulaValidator {
	// Release URLs use the standard validator (GitHub only)
	releaseValidator := NewURLValidator()

	// Homepage URLs can be any HTTPS URL (no host restriction)
	homepageValidator := &URLValidator{allowedHosts: []string{}}

	return &FormulaValidator{
		nameValidator:    NewPackageNameValidator(),
		versionValidator: NewVersionValidator(),
		urlValidator:     homepageValidator,
		releaseValidator: releaseValidator,
	}
}

// ValidateFile validates a formula YAML file
func (v *FormulaValidator) ValidateFile(path string) error {
	// Read file
	data, err := os.ReadFile(path) //nolint:gosec // G304: path from application config
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML
	var formula FormulaSchema
	if err := yaml.Unmarshal(data, &formula); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	// Validate structure
	return v.Validate(&formula)
}

// Validate validates a formula schema
func (v *FormulaValidator) Validate(formula *FormulaSchema) error {
	// Validate name
	if err := v.nameValidator.Validate(formula.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Validate version
	if err := v.versionValidator.Validate(formula.Version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	// Check required fields
	if formula.Description == "" {
		return errors.New("description is required")
	}

	if formula.Homepage == "" {
		return errors.New("homepage is required")
	}

	// Validate homepage URL
	if err := v.urlValidator.Validate(formula.Homepage); err != nil {
		return fmt.Errorf("invalid homepage: %w", err)
	}

	// Validate license
	if formula.License == "" {
		return errors.New("license is required")
	}

	// Validate releases
	if len(formula.Releases) == 0 {
		return errors.New("at least one release is required")
	}

	// Validate each release URL
	for platform, url := range formula.Releases {
		if platform == "" {
			return errors.New("release platform cannot be empty")
		}
		if err := v.releaseValidator.Validate(url); err != nil {
			return fmt.Errorf("invalid release URL for %s: %w", platform, err)
		}
	}

	return nil
}
