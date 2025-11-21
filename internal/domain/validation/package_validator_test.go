package validation

import "testing"

func TestValidate(t *testing.T) {
	v := NewPackageNameValidator()

	tests := []struct {
		name string
		want bool
	}{
		{"nano", true},
		{"microsoft-edge", true},
		{"", false},
		{"Visual-Studio", false},
		{"nano--editor", false},
	}

	for _, tt := range tests {
		err := v.Validate(tt.name)
		got := (err == nil)
		if got != tt.want {
			t.Errorf("Validate(%q) = %v, want %v (err: %v)", tt.name, got, tt.want, err)
		}
	}
}
