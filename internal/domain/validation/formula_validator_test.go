package validation

import (
	"testing"
)

func TestFormulaValidator_Validate(t *testing.T) {
	validator := NewFormulaValidator()

	tests := []struct {
		name    string
		formula *FormulaSchema
		wantErr bool
	}{
		{
			name: "valid formula",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				Homepage:    "https://www.nano-editor.org",
				License:     "GPL-3.0",
				Tags:        []string{"editor", "cli"},
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			formula: &FormulaSchema{
				Name:        "Invalid-Name",
				Version:     "1.0.0",
				Description: "Test",
				Homepage:    "https://example.com",
				License:     "MIT",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/test-1.0.0/test-1.0.0-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid version",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2; rm -rf /",
				Description: "Test",
				Homepage:    "https://example.com",
				License:     "MIT",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "missing description",
			formula: &FormulaSchema{
				Name:     "nano",
				Version:  "8.2",
				Homepage: "https://www.nano-editor.org",
				License:  "GPL-3.0",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "missing homepage",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				License:     "GPL-3.0",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid homepage URL",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				Homepage:    "http://insecure.com",
				License:     "GPL-3.0",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "missing license",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				Homepage:    "https://www.nano-editor.org",
				Releases: map[string]string{
					"darwin-x86_64": "https://github.com/ochairo/potions/releases/download/nano-8.2/nano-8.2-darwin-x86_64.tar.gz",
				},
			},
			wantErr: true,
		},
		{
			name: "missing releases",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				Homepage:    "https://www.nano-editor.org",
				License:     "GPL-3.0",
				Releases:    map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "invalid release URL",
			formula: &FormulaSchema{
				Name:        "nano",
				Version:     "8.2",
				Description: "Text editor",
				Homepage:    "https://www.nano-editor.org",
				License:     "GPL-3.0",
				Releases: map[string]string{
					"darwin-x86_64": "http://insecure.com/file.tar.gz",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.formula)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormulaValidator_ValidateFile(t *testing.T) {
	validator := NewFormulaValidator()

	// Test with non-existent file
	err := validator.ValidateFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("ValidateFile() should fail for non-existent file")
	}

	// Test with invalid YAML would require creating a temp file
	// Skipping for now as it requires file I/O setup
}
