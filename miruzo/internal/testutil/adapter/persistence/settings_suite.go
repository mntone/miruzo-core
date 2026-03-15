package persistence

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type SettingsSuite SuiteBase[persist.SettingsRepository]

func (ste SettingsSuite) RunTestGetValueReturnsNotFound(t *testing.T) {
	t.Helper()

	_, err := ste.Repository.GetValue(ste.Context, "key2")
	assert.ErrorIs(t, "GetValue() error", err, persist.ErrNotFound)
}

func (ste SettingsSuite) RunTestUpdateAndGetValue(t *testing.T) {
	t.Helper()

	err := ste.Repository.UpdateValue(ste.Context, "key1", "value")
	assert.NilError(t, "[1st] UpdateValue() error", err)

	value, err := ste.Repository.GetValue(ste.Context, "key1")
	assert.NilError(t, "[1st] GetValue() error", err)
	assert.Equal(t, "[1st] GetValue()", value, "value")

	err = ste.Repository.UpdateValue(ste.Context, "key1", "update")
	assert.NilError(t, "[2nd] UpdateValue() error", err)

	value, err = ste.Repository.GetValue(ste.Context, "key1")
	assert.NilError(t, "[2nd] GetValue() error", err)
	assert.Equal(t, "[2nd] GetValue()", value, "update")
}

func (ste SettingsSuite) RunTestUpdateValueReturnsCheckViolation(t *testing.T) {
	t.Helper()

	err := ste.Repository.UpdateValue(ste.Context, "k", "short")
	assert.ErrorIs(t, "UpdateValue() error", err, persist.ErrCheckViolation)

	err = ste.Repository.UpdateValue(ste.Context, "very_long_key", "long")
	assert.ErrorIs(t, "UpdateValue() error", err, persist.ErrCheckViolation)
}
