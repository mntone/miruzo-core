package bind

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindIntQueryReturnsValidValue(t *testing.T) {
	got, err := BindIntQuery[int]("key", []string{"-42"})
	assert.Nil(t, "BindIntQuery() error", err)
	assert.Equal(t, "BindIntQuery()", got, -42)
}

func TestBindIntQueryReturnsError(t *testing.T) {
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
			name:         "TooSmall",
			values:       []string{"-129"},
			errorType:    "invalid",
			errorMessage: "must be an integer",
		},
		{
			name:         "TooLarge",
			values:       []string{"128"},
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
			_, err := BindIntQuery[int8]("key", tt.values)
			assert.NotNil(t, "BindIntQuery() error", err)
			assert.Equal(t, "BindIntQuery().Type", err.Type, tt.errorType)
			assert.Equal(t, "BindIntQuery().Path", err.Path, "query.key")
			assert.Equal(t, "BindIntQuery().Message", err.Message, tt.errorMessage)
		})
	}
}
