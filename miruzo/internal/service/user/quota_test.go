package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/service/user"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func assertQuota(
	t *testing.T,
	name string,
	gotQuota model.Quota,
	wantPeriod model.PeriodType,
	wantResetAt time.Time,
	wantRemaining model.QuotaInt,
	wantLimit model.QuotaInt,
) {
	t.Helper()

	assert.Equal(t, name+".Period", gotQuota.Period, wantPeriod)
	assert.EqualFn(t, name+".ResetAt", gotQuota.ResetAt, wantResetAt)
	assert.Equal(t, name+".Limit", gotQuota.Limit, wantLimit)
	assert.Equal(t, name+".Remaining", gotQuota.Remaining, wantRemaining)
}

func TestGetQuotaReturnsRemainingAndResetAt(t *testing.T) {
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	resolver := period.NewDailyResolver(offset)

	service, err := user.New(
		stub.NewStubUserRepository(3),
		clock.NewFixedProvider(current),
		resolver,
		10,
	)
	assert.NilError(t, "user.New() error", err)

	response, err := service.GetQuota(context.Background())
	assert.NilError(t, "GetQuota() error", err)
	assertQuota(
		t,
		"GetQuota().Love",
		response.Love,
		model.PeriodTypeDaily,
		resolver.PeriodEnd(current),
		7, 10,
	)
}

func TestGetQuotaUsesLimitWhenNoLoveUsed(t *testing.T) {
	current := time.Date(2026, 1, 9, 16, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	resolver := period.NewDailyResolver(offset)

	service, err := user.New(
		stub.NewStubUserRepository(0),
		clock.NewFixedProvider(current),
		resolver,
		8,
	)
	assert.NilError(t, "user.New() error", err)

	response, err := service.GetQuota(context.Background())
	assert.NilError(t, "GetQuota() error", err)
	assertQuota(
		t,
		"GetQuota().Love",
		response.Love,
		model.PeriodTypeDaily,
		resolver.PeriodEnd(current),
		8, 8,
	)
}

func TestGetQuotaClampsRemainingToZero(t *testing.T) {
	current := time.Date(2026, 1, 18, 20, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	resolver := period.NewDailyResolver(offset)

	service, err := user.New(
		stub.NewStubUserRepository(99),
		clock.NewFixedProvider(current),
		resolver,
		5,
	)
	assert.NilError(t, "user.New() error", err)

	response, err := service.GetQuota(context.Background())
	assert.NilError(t, "GetQuota() error", err)
	assertQuota(
		t,
		"GetQuota().Love",
		response.Love,
		model.PeriodTypeDaily,
		resolver.PeriodEnd(current),
		0, 5,
	)
}

func TestGetQuotaReturnsNotFoundWhenSingletonUserMissing(t *testing.T) {
	current := time.Date(2026, 1, 20, 9, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	resolver := period.NewDailyResolver(offset)
	repo := stub.NewStubUserRepository(0)
	repo.GetError = persist.ErrNoRows

	service, err := user.New(
		repo,
		clock.NewFixedProvider(current),
		resolver,
		8,
	)
	assert.NilError(t, "user.New() error", err)

	_, err = service.GetQuota(context.Background())
	assert.ErrorIs(t, "GetQuota() error", err, serviceerror.ErrNotFound)
}
