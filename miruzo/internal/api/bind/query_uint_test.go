package bind

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindUintQueryReturnsValidValue(t *testing.T) {
	got, err := BindUintQuery[uint]("key", []string{"42"})
	assert.Nil(t, "BindUintQuery() error", err)
	assert.Equal(t, "BindUintQuery()", got, 42)
}

func TestBindUintQueryReturnsError(t *testing.T) {
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
			values:       []string{"20", "40"},
			errorType:    "duplicate",
			errorMessage: "must not be specified multiple times",
		},
		{
			name:         "String",
			values:       []string{"abc"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
		{
			name:         "Negative",
			values:       []string{"-42"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
		{
			name:         "TooSmall",
			values:       []string{"-1"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
		{
			name:         "TooLarge",
			values:       []string{"256"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
		{
			name:         "FloatingPoint",
			values:       []string{"12.3"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BindUintQuery[uint8]("key", tt.values)
			assert.NotNil(t, "BindUintQuery() error", err)
			assert.Equal(t, "BindUintQuery().Type", err.Type, tt.errorType)
			assert.Equal(t, "BindUintQuery().Path", err.Path, "query.key")
			assert.Equal(t, "BindUintQuery().Message", err.Message, tt.errorMessage)
		})
	}
}
