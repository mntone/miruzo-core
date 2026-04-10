package contract_test

import (
	"testing"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestSettingsRepositoryGetValueReturnsNotFound(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			_, err := ops.Settings().GetValue(t.Context(), "key2")
			assert.ErrorIs(t, "GetValue() error", err, persist.ErrNotFound)
		})
	})
}

func TestSettingsRepositoryUpdateAndGetValue(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			err := ops.Settings().UpdateValue(t.Context(), "key1", "value")
			assert.NilError(t, "[1st] UpdateValue() error", err)

			value, err := ops.Settings().GetValue(t.Context(), "key1")
			assert.NilError(t, "[1st] GetValue() error", err)
			assert.Equal(t, "[1st] GetValue()", value, "value")

			err = ops.Settings().UpdateValue(t.Context(), "key1", "update")
			assert.NilError(t, "[2nd] UpdateValue() error", err)

			value, err = ops.Settings().GetValue(t.Context(), "key1")
			assert.NilError(t, "[2nd] GetValue() error", err)
			assert.Equal(t, "[2nd] GetValue()", value, "update")
		})
	})
}

func TestSettingsRepositoryUpdateValueReturnsCheckViolation(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "TooShort",
			key:   "k",
			value: "short",
		},
		{
			name:  "TooLong",
			key:   "very_long_key",
			value: "long",
		},
	}

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(tt.name, func(t *testing.T) {
					err := ops.Settings().UpdateValue(t.Context(), tt.key, tt.value)
					assert.ErrorIs(t, "UpdateValue() error", err, persist.ErrCheckViolation)
				})
			})
		}
	})
}
