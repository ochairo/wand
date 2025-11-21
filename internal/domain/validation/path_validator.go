package validation

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PathValidator validates file paths for security and correctness.
type PathValidator struct{}

// NewPathValidator creates a new path validator.
func NewPathValidator() *PathValidator {
	return &PathValidator{}
}

// Validate checks if a path is valid and safe.
func (v *PathValidator) Validate(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	cleaned := filepath.Clean(path)
	if strings.Contains(cleaned, "..") {
		return fmt.Errorf("path contains directory traversal")
	}

	if strings.HasPrefix(path, "~") && !strings.HasPrefix(path, "~/") {
		return fmt.Errorf("invalid home directory reference")
	}

	return nil
}

// Sanitize removes potentially dangerous path components.
func (v *PathValidator) Sanitize(path string) string {
	cleaned := filepath.Clean(path)
	cleaned = strings.ReplaceAll(cleaned, "..", "")
	return cleaned
}
