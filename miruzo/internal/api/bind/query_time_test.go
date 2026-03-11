package bind

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestBindTimeQueryReturnsValidValue(t *testing.T) {
	got, err := BindTimeQuery("cursor", []string{"2026-03-02T10:20:30.123456Z"})
	assert.Nil(t, "BindTimeQuery() error", err)
	assert.EqualFn(t, "BindTimeQuery()", got, time.Date(2026, 3, 2, 10, 20, 30, 123456000, time.UTC))
}

func TestBindTimeQueryReturnsError(t *testing.T) {
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
			values:       []string{"2026-03-02T10:20:30.123456Z", "2026-03-02T10:20:30.123456Z"},
			errorType:    "duplicate",
			errorMessage: "must not be specified multiple times",
		},
		{
			name:         "Invalid",
			values:       []string{"not-a-time"},
			errorType:    "invalid",
			errorMessage: "must be ISO8601 timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BindTimeQuery("cursor", tt.values)
			assert.NotNil(t, "BindTimeQuery() error", err)
			assert.Equal(t, "BindTimeQuery().Type", err.Type, tt.errorType)
			assert.Equal(t, "BindTimeQuery().Path", err.Path, "query.cursor")
			assert.Equal(t, "BindTimeQuery().Message", err.Message, tt.errorMessage)
		})
	}
}
