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

func TestLoveCancelRestoresPreviousLove(t *testing.T) {
	ingestID := model.IngestIDType(1)
	previous := mo.Some(time.Date(2026, 1, 1, 20, 0, 0, 0, time.UTC))
	current := time.Date(2026, 1, 2, 23, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(1, model.Stats{
		IngestID:     ingestID,
		Score:        120,
		FirstLovedAt: previous,
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

	response, err := service.LoveCancel(context.Background(), ingestID)
	assert.NilError(t, "LoveCancel() error", err)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 0)

	assert.Equal(t, "LoveCancel().Quota.Period", response.Quota.Period, model.PeriodTypeDaily)
	assert.EqualFn(t, "LoveCancel().Quota.ResetAt", response.Quota.ResetAt, resolver.PeriodEnd(current))
	assert.Equal(t, "LoveCancel().Quota.Limit", response.Quota.Limit, 3)
	assert.Equal(t, "LoveCancel().Quota.Remaining", response.Quota.Remaining, 3)

	assert.Equal(t, "LoveCancel().Stats.Score", response.Stats.Score, 120+scoreCalc.LoveCanceledDelta())
	assert.EqualFn(t, "LoveCancel().Stats.FirstLovedAt", response.Stats.FirstLovedAt, previous)
	assert.EqualFn(t, "LoveCancel().Stats.LastLovedAt", response.Stats.LastLovedAt, previous)
}

func TestLoveCancelRestoresNullWhenNoPreviousLove(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := mo.Some(time.Date(2026, 1, 3, 19, 0, 0, 0, time.UTC))
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(2, model.Stats{
		IngestID:     ingestID,
		Score:        120,
		FirstLovedAt: current,
		LastLovedAt:  current,
	})
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current.MustGet()),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	response, err := service.LoveCancel(context.Background(), ingestID)
	assert.NilError(t, "LoveCancel() error", err)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 1)

	assert.Equal(t, "LoveCancel().Quota.Period", response.Quota.Period, model.PeriodTypeDaily)
	assert.EqualFn(t, "LoveCancel().Quota.ResetAt", response.Quota.ResetAt, resolver.PeriodEnd(current.MustGet()))
	assert.Equal(t, "LoveCancel().Quota.Limit", response.Quota.Limit, 3)
	assert.Equal(t, "LoveCancel().Quota.Remaining", response.Quota.Remaining, 2)

	assert.Equal(t, "LoveCancel().Stats.Score", response.Stats.Score, 120+scoreCalc.LoveCanceledDelta())
	assert.IsAbsent(t, "LoveCancel().Stats.FirstLovedAt", response.Stats.FirstLovedAt)
	assert.IsAbsent(t, "LoveCancel().Stats.LastLovedAt", response.Stats.LastLovedAt)
}

func TestLoveCancelReturnsConflictWhenNoLoveInPeriod(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := time.Date(2026, 1, 14, 17, 0, 0, 0, time.UTC)
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(1, model.Stats{
		IngestID: ingestID,
		Score:    120,
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

	_, err = service.LoveCancel(context.Background(), ingestID)
	assert.ErrorIs(t, "LoveCancel() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 1)
	assert.Empty(t, "action count", manager.Action.Store)
}

func TestLoveCancelRollsBackWhenActionCreateFails(t *testing.T) {
	ingestID := model.IngestIDType(1)
	current := mo.Some(time.Date(2026, 1, 17, 14, 0, 0, 0, time.UTC))
	offset := 5 * time.Hour
	manager := stub.NewStubPersistenceManager(2, model.Stats{
		IngestID:     ingestID,
		Score:        100,
		FirstLovedAt: current,
		LastLovedAt:  current,
	})
	manager.Action.CreateError = persist.ErrUnavailable
	resolver := period.NewDailyResolver(offset)
	scoreCalc := testutilDomain.NewTestScoreCalculator(resolver)

	service, err := reaction.New(
		manager,
		clock.NewFixedProvider(current.MustGet()),
		resolver,
		scoreCalc,
		3,
		180,
	)
	assert.NilError(t, "reaction.New() error", err)

	_, err = service.LoveCancel(context.Background(), ingestID)
	assert.ErrorIs(t, "LoveCancel() error", err, serviceerror.ErrServiceUnavailable)
	assert.Equal(t, "daily_love_used", manager.User.DailyLoveUsed, 2)

	stats := manager.Stats.Store[ingestID]
	assert.Equal(t, "stats.Score", stats.Score, 100)
	assert.Equal(t, "stats.FirstLovedAt", stats.FirstLovedAt, current)
	assert.Equal(t, "stats.LastLovedAt", stats.LastLovedAt, current)
	assert.Empty(t, "action count", manager.Action.Store)
}
