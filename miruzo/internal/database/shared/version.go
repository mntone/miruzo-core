package shared

import "fmt"

type Version struct {
	Major, Minor, Patch int
}

func (v Version) LessThan(min Version) bool {
	if v.Major != min.Major {
		return v.Major < min.Major
	}
	if v.Minor != min.Minor {
		return v.Minor < min.Minor
	}
	return v.Patch < min.Patch
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ParseVersion parses the leading semantic-like "major.minor.patch" numbers.
// Any trailing suffix after the patch number is ignored by design.
func ParseVersion(s string) (Version, error) {
	var v Version
	if _, err := fmt.Sscanf(s, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch); err != nil {
		return Version{}, fmt.Errorf("invalid version %q: %w", s, err)
	}
	return v, nil
}
