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
			assert.ErrorIs(t, "GetValue() error", err, persist.ErrNoRows)
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
			name:  "TooShortKey",
			key:   "k",
			value: "short",
		},
		{
			name:  "TooLongKey",
			key:   "very_long",
			value: "long",
		},
		{
			name:  "InvalidKeyChars",
			key:   "dash-key",
			value: "invalid chars: dash",
		},
		{
			name:  "InvalidKeyUppercase",
			key:   "UPPER123",
			value: "invalid chars: uppercase",
		},
		{
			name:  "TooLongValue",
			key:   "long_val",
			value: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
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
