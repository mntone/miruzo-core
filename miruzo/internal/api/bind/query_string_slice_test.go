package bind

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindStringSliceQueryReturnsValidValues(t *testing.T) {
	got, err := BindStringSliceQuery("key", []string{"jpeg+webp+png"}, "+")
	assert.Nil(t, "BindStringSliceQuery() error", err)
	assert.LenIs(t, "BindStringSliceQuery()", got, 3)
	assert.Equal(t, "BindStringSliceQuery()[0]", got[0], "jpeg")
	assert.Equal(t, "BindStringSliceQuery()[1]", got[1], "webp")
	assert.Equal(t, "BindStringSliceQuery()[2]", got[2], "png")
}

func TestBindStringSliceQuerySkipsEmptyItems(t *testing.T) {
	got, err := BindStringSliceQuery("key", []string{"jpeg++webp+"}, "+")
	assert.Nil(t, "BindStringSliceQuery() error", err)
	assert.LenIs(t, "BindStringSliceQuery()", got, 2)
	assert.Equal(t, "BindStringSliceQuery()[0]", got[0], "jpeg")
	assert.Equal(t, "BindStringSliceQuery()[1]", got[1], "webp")
}

func TestBindStringSliceQueryReturnsError(t *testing.T) {
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
			values:       []string{"jpeg+webp", "png"},
			errorType:    "duplicate",
			errorMessage: "must not be specified multiple times",
		},
		{
			name:         "OnlySeparators",
			values:       []string{"+++"},
			errorType:    "invalid",
			errorMessage: "must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BindStringSliceQuery("key", tt.values, "+")
			assert.NotNil(t, "BindStringSliceQuery() error", err)
			assert.Equal(t, "BindStringSliceQuery().Type", err.Type, tt.errorType)
			assert.Equal(t, "BindStringSliceQuery().Path", err.Path, "query.key")
			assert.Equal(t, "BindStringSliceQuery().Message", err.Message, tt.errorMessage)
		})
	}
}
