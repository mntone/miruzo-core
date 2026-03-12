package item

import (
	"net/url"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindLevelQueryReturnsValidValue(t *testing.T) {
	rich, err := bindLevelQuery("level", []string{"rich"})
	assert.Nil(t, "bindLevelQuery() error", err)
	assert.Equal(t, "bindLevelQuery()", rich, true)

	rich, err = bindLevelQuery("level", []string{"default"})
	assert.Nil(t, "bindLevelQuery() error", err)
	assert.Equal(t, "bindLevelQuery()", rich, false)
}

func TestBindLevelQueryReturnsError(t *testing.T) {
	tests := []struct {
		name         string
		values       []string
		errorType    string
		errorMessage string
	}{
		{
			name:         "Empty",
			values:       []string{},
			errorType:    "invalid",
			errorMessage: "must not be empty",
		},
		{
			name:         "Duplicate",
			values:       []string{"rich", "default"},
			errorType:    "duplicate",
			errorMessage: "must not be specified multiple times",
		},
		{
			name:         "Invalid",
			values:       []string{"full"},
			errorType:    "invalid",
			errorMessage: "must be \"default\" or \"rich\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bindLevelQuery("level", tt.values)
			assert.NotNil(t, "bindLevelQuery() error", err)
			assert.Equal(t, "bindLevelQuery().Type", err.Type, tt.errorType)
			assert.Equal(t, "bindLevelQuery().Path", err.Path, "query.level")
			assert.Equal(t, "bindLevelQuery().Message", err.Message, tt.errorMessage)
		})
	}
}

func TestBindParamsReturnsErrorForUnsupportedKey(t *testing.T) {
	_, errs := bindParams(url.Values{
		"unknown": []string{"1"},
	})
	assert.LenIs(t, "bindParams() error", errs, 1)
	assert.Equal(t, "err[0].Type", errs[0].Type, "unsupported")
	assert.Equal(t, "err[0].Path", errs[0].Path, "query.unknown")
	assert.Equal(t, "err[0].Message", errs[0].Message, "is not supported")
}
