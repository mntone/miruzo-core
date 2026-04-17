package mysql

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/database/shared"
)

func TestSupportsMySQLVersionAcceptsBoundary(t *testing.T) {
	version := shared.Version{
		Major: 8,
		Minor: 0,
		Patch: 16,
	}
	if !supportsMySQLCheckVersion(version) {
		t.Fatalf("supportsMySQLCheckVersion(%q) = false, want true", "8.0.16")
	}
}

func TestSupportsMySQLVersionRejectsOlderVersion(t *testing.T) {
	version := shared.Version{
		Major: 8,
		Minor: 0,
		Patch: 15,
	}
	if supportsMySQLCheckVersion(version) {
		t.Fatalf("supportsMySQLCheckVersion(%q) = true, want false", "8.0.15")
	}
}
