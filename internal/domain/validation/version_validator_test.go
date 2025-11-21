package validation

import "testing"

func TestVersionValidator_Validate(t *testing.T) {
	v := NewVersionValidator()

	tests := []struct {
		version string
		want    bool
	}{
		{"8.7", true},
		{"8.7.0", true},
		{"1.2.3", true},
		{"1.0.0-beta", true},
		{"", false},
		{"v8.7", false},
		{"8.7.0.1", false},
		{"abc", false},
	}

	for _, tt := range tests {
		err := v.Validate(tt.version)
		got := (err == nil)
		if got != tt.want {
			t.Errorf("Validate(%q) = %v, want %v", tt.version, got, tt.want)
		}
	}
}
