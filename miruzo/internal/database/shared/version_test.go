package shared

import (
	"strings"
	"testing"
)

func TestParseVersion(t *testing.T) {
	testCases := []struct {
		version  string
		expected Version
	}{
		{
			version: "3.35.0",
			expected: Version{
				Major: 3,
				Minor: 35,
				Patch: 0,
			},
		},
		{
			version: "3.45.0",
			expected: Version{
				Major: 3,
				Minor: 45,
				Patch: 0,
			},
		},
		{
			version: "8.0.16foo",
			expected: Version{
				Major: 8,
				Minor: 0,
				Patch: 16,
			},
		},
		{
			version: "8.0.36-0ubuntu0.22.04.1",
			expected: Version{
				Major: 8,
				Minor: 0,
				Patch: 36,
			},
		},
		{
			version: "9.6.0",
			expected: Version{
				Major: 9,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			version: "10.4.7-MariaDB",
			expected: Version{
				Major: 10,
				Minor: 4,
				Patch: 7,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.version, func(t *testing.T) {
			got, err := ParseVersion(testCase.version)
			if err != nil {
				t.Fatalf("ParseVersion(%q) returned unexpected error: %v", testCase.version, err)
			}

			if got != testCase.expected {
				t.Fatalf("ParseVersion(%q) = %+v, want %+v", testCase.version, got, testCase.expected)
			}
		})
	}
}

func TestParseVersionReturnsErrorForInvalidFormat(t *testing.T) {
	testCases := []string{
		"3",
		"3.35",
		"x.35.0",
		"3.35.beta",
	}

	for _, version := range testCases {
		t.Run(version, func(t *testing.T) {
			_, err := ParseVersion(version)
			if err == nil {
				t.Fatalf("ParseVersion(%q) expected error but got nil", version)
			}

			if !strings.Contains(err.Error(), "invalid version") {
				t.Fatalf("ParseVersion(%q) error = %q, want to include %q", version, err.Error(), "invalid version")
			}
		})
	}
}
