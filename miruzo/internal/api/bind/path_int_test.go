package bind_test

import (
	"net/http/httptest"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/api/apierror"
	"github.com/mntone/miruzo-core/miruzo/internal/api/bind"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindIntPathReturnsValidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/123", nil)
	req.SetPathValue("id", "123")

	got, err := bind.BindIntPath[int](req, "id")
	assert.Nil(t, "BindIntPath() error", err)
	assert.Equal(t, "BindIntPath()", got, 123)
}

func TestBindIntPathReturnsInvalidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/abc", nil)
	req.SetPathValue("id", "abc")

	_, err := bind.BindIntPath[int](req, "id")
	assert.NotNil(t, "BindIntPath() error", err)
	assert.Equal(t, "BindIntPath().Path", err.Path, "path.id")
	assert.Equal(t, "BindIntPath().Type", err.Type, apierror.FieldErrorTypeInvalid)
	assert.Equal(t, "BindIntPath().Message", err.Message, "must be an integer")
}

func TestParseIntPathIsMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "/images/", nil)

	_, err := bind.BindIntPath[int](req, "id")
	assert.NotNil(t, "BindIntPath() error", err)
	assert.Equal(t, "BindIntPath().Path", err.Path, "path.id")
	assert.Equal(t, "BindIntPath().Type", err.Type, apierror.FieldErrorTypeRequired)
	assert.Equal(t, "BindIntPath().Message", err.Message, "is required")
}
