package persistence

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type UserSuite SuiteBase[persist.UserRepository]

func (ste UserSuite) RunTestGetSingletonUser(t *testing.T) {
	t.Helper()

	user, err := ste.Repository.GetSingletonUser(ste.Context)
	assert.NilError(t, "GetSingletonUser() error", err)
	assert.Equal(t, "user.ID", user.ID, 1)
	assert.Equal(t, "user.DailyLoveUsed", user.DailyLoveUsed, 0)
}

// --- increment daily_love_used ---

func (ste UserSuite) RunTestIncrementDailyLoveUsedIncrements(t *testing.T) {
	t.Helper()

	dailyLoveUsed, err := ste.Repository.IncrementDailyLoveUsed(ste.Context, 2)
	assert.NilError(t, "IncrementDailyLoveUsed() error", err)
	assert.Equal(t, "IncrementDailyLoveUsed()", dailyLoveUsed, 1)
}

func (ste UserSuite) RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenMissing(t *testing.T) {
	t.Helper()

	ste.Operations.MustRemoveUser(t)

	_, err := ste.Repository.IncrementDailyLoveUsed(ste.Context, 2)
	assert.ErrorIs(t, "IncrementDailyLoveUsed() error", err, persist.ErrQuotaExceeded)
}

func (ste UserSuite) RunTestIncrementDailyLoveUsedReturnsQuotaExceededWhenLimitReached(t *testing.T) {
	t.Helper()

	ste.Operations.MustSetDailyLoveUsed(t, 2)

	_, err := ste.Repository.IncrementDailyLoveUsed(ste.Context, 2)
	assert.ErrorIs(t, "IncrementDailyLoveUsed() error", err, persist.ErrQuotaExceeded)
}
