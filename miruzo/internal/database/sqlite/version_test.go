package sqlite

import (
	"strings"
	"testing"
)

func TestParseSQLiteVersion(t *testing.T) {
	testCases := []struct {
		version  string
		expected sqliteVersion
	}{
		{
			version: "3.35.0",
			expected: sqliteVersion{
				major: 3,
				minor: 35,
				patch: 0,
			},
		},
		{
			version: "3.45.0",
			expected: sqliteVersion{
				major: 3,
				minor: 45,
				patch: 0,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.version, func(t *testing.T) {
			got, err := parseSQLiteVersion(testCase.version)
			if err != nil {
				t.Fatalf("parseSQLiteVersion(%q) returned unexpected error: %v", testCase.version, err)
			}

			if got != testCase.expected {
				t.Fatalf("parseSQLiteVersion(%q) = %+v, want %+v", testCase.version, got, testCase.expected)
			}
		})
	}
}

func TestParseSQLiteVersionReturnsErrorForInvalidFormat(t *testing.T) {
	testCases := []string{
		"3",
		"3.35",
		"x.35.0",
		"3.35.beta",
	}

	for _, version := range testCases {
		t.Run(version, func(t *testing.T) {
			_, err := parseSQLiteVersion(version)
			if err == nil {
				t.Fatalf("parseSQLiteVersion(%q) expected error but got nil", version)
			}

			if !strings.Contains(err.Error(), "invalid sqlite_version") {
				t.Fatalf("parseSQLiteVersion(%q) error = %q, want to include %q", version, err.Error(), "invalid sqlite_version")
			}
		})
	}
}

func TestSupportsSQLiteReturningVersionAcceptsBoundary(t *testing.T) {
	if !supportsSQLiteReturningVersion(sqliteVersion{3, 35, 0}) {
		t.Fatalf("supportsSQLiteReturningVersion(%q) = false, want true", "3.35.0")
	}
}

func TestSupportsSQLiteReturningVersionRejectsOlderVersion(t *testing.T) {
	if supportsSQLiteReturningVersion(sqliteVersion{3, 34, 1}) {
		t.Fatalf("supportsSQLiteReturningVersion(%q) = true, want false", "3.34.1")
	}
}
