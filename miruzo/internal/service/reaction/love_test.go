package reaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	testutilDomain "github.com/mntone/miruzo-core/miruzo/internal/testutil/domain"
	"github.com/samber/mo"
)

func TestLoveUpdates(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(0, model.Stats{
		IngestID: ingestID,
		Score:    100,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	response, err := service.Love(context.Background(), ingestID)
	assert.NilError(t, "Love() error", err)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 1)

	assert.Equal(t, "Love().Quota.Period", response.Quota.Period, model.PeriodTypeDaily)
	assert.EqualFn(t, "Love().Quota.ResetAt", response.Quota.ResetAt, resolver.PeriodEnd(current))
	assert.Equal(t, "Love().Quota.Limit", response.Quota.Limit, 3)
	assert.Equal(t, "Love().Quota.Remaining", response.Quota.Remaining, 2)

	assert.Equal(t, "Love().Stats.Score", response.Stats.Score, 100+scoreCalc.LoveDelta())
	assert.EqualFn(t, "Love().Stats.FirstLovedAt", response.Stats.FirstLovedAt, mo.Some(current))
	assert.EqualFn(t, "Love().Stats.LastLovedAt", response.Stats.LastLovedAt, mo.Some(current))
}

func TestLoveReturnsConflictWhenAlreadyLovedToday(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 5, 19, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(1, model.Stats{
		IngestID:     ingestID,
		Score:        120,
		FirstLovedAt: mo.Some(current),
		LastLovedAt:  mo.Some(current),
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.Love(context.Background(), ingestID)
	assert.ErrorIs(t, "Love() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 1)
	assert.Empty(t, "action count", manager.Action.Store)
}

func TestLoveReturnsTooManyRequestsWhenQuotaExceeded(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 12, 22, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(1, model.Stats{
		IngestID: ingestID,
		Score:    100,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		1,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.Love(context.Background(), ingestID)
	assert.ErrorIs(t, "Love() error", err, serviceerror.ErrTooManyRequests)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 1)
	assert.Empty(t, "action count", manager.Action.Store)
}

func TestLoveRollsBackWhenActionCreateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 18, 10, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(0, model.Stats{
		IngestID: ingestID,
		Score:    100,
	})
	manager.Action.CreateError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.Love(context.Background(), ingestID)
	assert.ErrorIs(t, "Love() error", err, serviceerror.ErrServiceUnavailable)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 0)

	stats := manager.Stats.Store[ingestID]
	assert.Equal(t, "stats.Score", stats.Score, 100)
	assert.IsAbsent(t, "stats.FirstLovedAt", stats.FirstLovedAt)
	assert.IsAbsent(t, "stats.LastLovedAt", stats.LastLovedAt)
	assert.Empty(t, "action count", manager.Action.Store)
}
