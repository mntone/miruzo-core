package sqlite

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

func TestSupportsSQLiteReturningAndStrictVersionAcceptsBoundary(t *testing.T) {
	version := shared.Version{
		Major: 3,
		Minor: 37,
		Patch: 0,
	}
	if !supportsSQLiteReturningAndStrictVersion(version) {
		t.Fatalf("supportsSQLiteReturningAndStrictVersion(%q) = false, want true", "3.37.0")
	}
}

func TestSupportsSQLiteReturningAndStrictVersionRejectsOlderVersion(t *testing.T) {
	version := shared.Version{
		Major: 3,
		Minor: 34,
		Patch: 1,
	}
	if supportsSQLiteReturningAndStrictVersion(version) {
		t.Fatalf("supportsSQLiteReturningAndStrictVersion(%q) = true, want false", "3.34.1")
	}
}
