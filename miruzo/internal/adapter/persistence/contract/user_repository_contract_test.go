package contract_test

import (
	"fmt"
	"testing"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

// --- Schema ---

func TestUserRepositorySchemaRejectsInvalidDailyLoveUsed(t *testing.T) {
	tests := []int32{-1, 101}

	stmt := "UPDATE users SET daily_love_used=%s WHERE id=1"

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
				t.Run(fmt.Sprintf("daily_love_used=%d", tt), func(t *testing.T) {
					ops.AssertExecErrorIs(
						t,
						c.DBErrorMappingDefault,
						persist.ErrCheckViolation,
						fmt.Sprintf(stmt, ops.Param(1)),
						tt,
					)
				})
			})
		}
	})
}

// --- Get ---

func TestUserRepositoryGetReturnsUser(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			user, err := ops.User().Get(t.Context())
			assert.NilError(t, "Get() error", err)
			assert.Equal(t, "user.ID", user.ID, 1)
			assert.Equal(t, "user.DailyLoveUsed", user.DailyLoveUsed, 0)
		})
	})
}

func TestUserRepositoryGetReturnsNotFoundWhenUserMissing(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustRemoveUser(t)

			_, err := ops.User().Get(t.Context())
			assert.ErrorIs(t, "Get() error", err, persist.ErrNotFound)
		})
	})
}

// --- IncrementDailyLoveUsed ---

func TestUserRepositoryIncrementDailyLoveUsedIncrements(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			dailyLoveUsed, err := ops.User().IncrementDailyLoveUsed(t.Context(), 2)
			assert.NilError(t, "IncrementDailyLoveUsed() error", err)
			assert.Equal(t, "IncrementDailyLoveUsed()", dailyLoveUsed, 1)
		})
	})
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenUserMissing(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustRemoveUser(t)

			_, err := ops.User().IncrementDailyLoveUsed(t.Context(), 2)
			assert.ErrorIs(t, "IncrementDailyLoveUsed() error", err, persist.ErrQuotaExceeded)
		})
	})
}

func TestUserRepositoryIncrementDailyLoveUsedReturnsQuotaExceededWhenAtLimit(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustSetDailyLoveUsed(t, 2)

			_, err := ops.User().IncrementDailyLoveUsed(t.Context(), 2)
			assert.ErrorIs(t, "IncrementDailyLoveUsed() error", err, persist.ErrQuotaExceeded)
		})
	})
}

// --- DecrementDailyLoveUsed ---

func TestUserRepositoryDecrementDailyLoveUsedDecrements(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustSetDailyLoveUsed(t, 1)

			dailyLoveUsed, err := ops.User().DecrementDailyLoveUsed(t.Context())
			assert.NilError(t, "DecrementDailyLoveUsed() error", err)
			assert.Equal(t, "DecrementDailyLoveUsed()", dailyLoveUsed, 0)
		})
	})
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsNotFound(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustRemoveUser(t)

			_, err := ops.User().DecrementDailyLoveUsed(t.Context())
			assert.ErrorIs(t, "DecrementDailyLoveUsed() error", err, persist.ErrNotFound)
		})
	})
}

func TestUserRepositoryDecrementDailyLoveUsedReturnsQuotaUnderflow(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			_, err := ops.User().DecrementDailyLoveUsed(t.Context())
			assert.ErrorIs(t, "DecrementDailyLoveUsed() error", err, persist.ErrQuotaUnderflow)
		})
	})
}

// --- ResetDailyLoveUsed ---

func TestUserRepositoryResetDailyLoveUsedResets(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustSetDailyLoveUsed(t, 5)

			err := ops.User().ResetDailyLoveUsed(t.Context())
			assert.NilError(t, "ResetDailyLoveUsed() error", err)

			user, err := ops.User().Get(t.Context())
			assert.NilError(t, "Get() error", err)
			assert.Equal(t, "Get().DailyLoveUsed", user.DailyLoveUsed, 0)
		})
	})
}

func TestUserRepositoryResetDailyLoveUsedReturnsNotFoundWhenUserMissing(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustRemoveUser(t)

			err := ops.User().ResetDailyLoveUsed(t.Context())
			assert.ErrorIs(t, "ResetDailyLoveUsed() error", err, persist.ErrNotFound)
		})
	})
}
