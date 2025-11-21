package types

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string
	Build string
}

// ParseVersion creates a Version from a version string
func ParseVersion(v string) (*Version, error) {
	if v == "" {
		return nil, fmt.Errorf("version string is empty")
	}

	v = strings.TrimPrefix(v, "v")

	var build string
	if idx := strings.Index(v, "+"); idx != -1 {
		build = v[idx+1:]
		v = v[:idx]
	}

	var pre string
	if idx := strings.Index(v, "-"); idx != -1 {
		pre = v[idx+1:]
		v = v[:idx]
	}

	parts := strings.Split(v, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	version := &Version{
		Pre:   pre,
		Build: build,
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}
	version.Major = major

	if len(parts) >= 2 {
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", parts[1])
		}
		version.Minor = minor
	}

	if len(parts) >= 3 {
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", parts[2])
		}
		version.Patch = patch
	}

	return version, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}

// Compare compares two versions
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	if v.Pre == "" && other.Pre != "" {
		return 1
	}
	if v.Pre != "" && other.Pre == "" {
		return -1
	}
	if v.Pre != other.Pre {
		if v.Pre < other.Pre {
			return -1
		}
		return 1
	}

	return 0
}

// LessThan returns true if v < other
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// GreaterThan returns true if v > other
func (v *Version) GreaterThan(other *Version) bool {
	return v.Compare(other) > 0
}

// Equal returns true if v == other
func (v *Version) Equal(other *Version) bool {
	return v.Compare(other) == 0
}
