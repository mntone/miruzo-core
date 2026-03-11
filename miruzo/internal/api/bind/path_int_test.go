package bind_test

import (
	"net/http/httptest"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestParseIntPathReturnsValidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/123", nil)
	req.SetPathValue("id", "123")

	got, err := bind.ParseIntPath[int](req, "id")
	if err != nil {
		t.Fatalf("ParseIntPath() error = \"%v\", want nil", err)
	}
	assert.Equal(t, "ParseIntPath()", got, 123)
}

func TestParseIntPathReturnsInvalidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/abc", nil)
	req.SetPathValue("id", "abc")

	_, err := bind.ParseIntPath[int](req, "id")
	if err == nil {
		t.Fatalf("ParseIntPath() error = nil, want non-empty")
	}
	assert.Equal(t, "err[0].Path", err[0].Path, "path.id")
	assert.Equal(t, "err[0].Type", err[0].Type, "invalid_type")
	assert.Equal(t, "err[0].Message", err[0].Message, "id must be an integer")
}

func TestParseIntPathIsMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/", nil)

	_, err := bind.ParseIntPath[int](req, "id")
	if err == nil {
		t.Fatalf("ParseIntPath() error = nil, want non-empty")
	}
	assert.Equal(t, "err[0].Path", err[0].Path, "path.id")
	assert.Equal(t, "err[0].Type", err[0].Type, "missing")
	assert.Equal(t, "err[0].Message", err[0].Message, "id is required")
}
