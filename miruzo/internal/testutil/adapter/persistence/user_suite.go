package persistence

import (
	"fmt"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type UserSuite SuiteBase[persist.UserRepository]

func (ste UserSuite) RunTestUserSchemaRejectsInvalidDailyLoveUsed(t *testing.T) {
	t.Helper()

	tests := []struct {
		name          string
		dailyLoveUsed int32
		wantErr       error
	}{
		{
			name:          "daily_love_used=-1",
			dailyLoveUsed: -1,
			wantErr:       persist.ErrCheckViolation,
		},
		{
			name:          "daily_love_used=101",
			dailyLoveUsed: 101,
			wantErr:       persist.ErrCheckViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"UPDATE users SET daily_love_used=%d WHERE id=1",
				tt.dailyLoveUsed,
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "update error", err, tt.wantErr)
		})
	}
}

func (ste UserSuite) RunTestGetReturnsUser(t *testing.T) {
	t.Helper()

	user, err := ste.Repository.Get(ste.Context)
	assert.NilError(t, "Get() error", err)
	assert.Equal(t, "user.ID", user.ID, 1)
	assert.Equal(t, "user.DailyLoveUsed", user.DailyLoveUsed, 0)
}

func (ste UserSuite) RunTestGetReturnsNotFoundWhenMissing(t *testing.T) {
	t.Helper()

	ste.Operations.MustRemoveUser(t)

	_, err := ste.Repository.Get(ste.Context)
	assert.ErrorIs(t, "Get() error", err, persist.ErrNotFound)
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

// --- decrement daily_love_used ---

func (ste UserSuite) RunTestDecrementDailyLoveUsedDecrements(t *testing.T) {
	t.Helper()

	ste.Operations.MustSetDailyLoveUsed(t, 1)

	dailyLoveUsed, err := ste.Repository.DecrementDailyLoveUsed(ste.Context)
	assert.NilError(t, "DecrementDailyLoveUsed() error", err)
	assert.Equal(t, "DecrementDailyLoveUsed()", dailyLoveUsed, 0)
}

func (ste UserSuite) RunTestDecrementDailyLoveUsedReturnsNotFound(t *testing.T) {
	t.Helper()

	ste.Operations.MustRemoveUser(t)

	_, err := ste.Repository.DecrementDailyLoveUsed(ste.Context)
	assert.ErrorIs(t, "DecrementDailyLoveUsed() error", err, persist.ErrNotFound)
}

func (ste UserSuite) RunTestDecrementDailyLoveUsedReturnsQuotaUnderflow(t *testing.T) {
	t.Helper()

	_, err := ste.Repository.DecrementDailyLoveUsed(ste.Context)
	assert.ErrorIs(t, "DecrementDailyLoveUsed() error", err, persist.ErrQuotaUnderflow)
}
