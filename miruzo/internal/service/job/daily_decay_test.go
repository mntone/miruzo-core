package job_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/job"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/domain"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

func newDailyDecayServiceFixture(
	dailyLoveUsed int32,
	evaluatedAt time.Time,
	statsEntries ...model.Stats,
) (service *job.DailyDecayService, mgr stub.PersistenceManager) {
	mgr = stub.NewStubPersistenceManager(dailyLoveUsed, statsEntries...)
	resolver := period.NewDailyResolver(5 * time.Hour)
	service = job.NewDailyDecay(
		mgr,
		clock.NewFixedProvider(evaluatedAt),
		resolver,
		domain.NewTestScoreCalculator(resolver),
	)
	return
}

func TestDailyDecayServiceApplyUpdatesScoresAndResetsDailyLoveUsed(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		3,
		baseTime.Add(1*time.Second),

		// No Viewed
		mb.Stats(1).Score(80).Build(),

		// Viewed (Normal Score)
		mb.Stats(2).
			Score(120).
			Viewed(24, baseTime.Add(-30*time.Minute)).
			Build(),

		// High Score
		mb.Stats(3).
			Score(180).
			Viewed(36, baseTime.Add(-2*time.Hour)).
			Build(),

		// Has daily decay action
		mb.Stats(4).
			Score(138).
			Viewed(52, baseTime.Add(-4*time.Hour)).
			EvaluateScore(baseTime).
			Build(),
	)
	evaluatedAt := baseTime.Add(1 * time.Second)
	_, _ = mgr.Action.Create(context.Background(), 4, model.ActionTypeDecay, baseTime)

	err := service.ApplyDailyDecay(context.Background())
	assert.NilError(t, "ApplyDailyDecay() error", err)
	assert.Equal(t, "IterateStatsForDailyDecay() batchCount", mgr.Stats.IterateStatsForDailyDecayArgs[0], 500)

	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.Stats.Store[1].Score, 80)
	assert.IsAbsent(t, "store.Stats(IngestID=1).ScoreEvaluatedAt", mgr.Stats.Store[1].ScoreEvaluatedAt)

	assert.Equal(t, "store.Stats(IngestID=2).Score", mgr.Stats.Store[2].Score, 120-2)
	assert.EqualFn(t, "store.Stats(IngestID=2).ScoreEvaluatedAt", mgr.Stats.Store[2].ScoreEvaluatedAt, mo.Some(evaluatedAt))

	assert.Equal(t, "store.Stats(IngestID=3).Score", mgr.Stats.Store[3].Score, 180-3)
	assert.EqualFn(t, "store.Stats(IngestID=3).ScoreEvaluatedAt", mgr.Stats.Store[3].ScoreEvaluatedAt, mo.Some(evaluatedAt))

	assert.Equal(t, "store.Stats(IngestID=4).Score", mgr.Stats.Store[4].Score, 138)
	assert.EqualFn(t, "store.Stats(IngestID=4).ScoreEvaluatedAt", mgr.Stats.Store[4].ScoreEvaluatedAt, mo.Some(baseTime))

	assert.LenIs(t, "store.Action count", mgr.Action.Store, 3)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, 0)
}

func TestDailyDecayServiceApplyReturnsIterateError(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	expectedErr := errors.New("iterate failed")
	mgr.Stats.IterateStatsForDailyDecayError = expectedErr

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, expectedErr)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsExistsSinceError(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.Action.ExistsSinceError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.Stats.Store[1].Score, model.ScoreType(120))
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsCreateError(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.Action.CreateError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Empty(t, "store.Action", mgr.Action.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsApplyDecayError(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.Stats.ApplyDecayError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.Stats.Store[1].Score, model.ScoreType(120))
	assert.Empty(t, "store.Action", mgr.Action.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsResetDailyLoveUsedError(t *testing.T) {
	baseTime := mb.GetDefaultStatsBaseTime()
	service, mgr := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.User.ResetError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.Stats.Store[1].Score, model.ScoreType(120))
	assert.Empty(t, "store.Action", mgr.Action.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.User.DailyLoveUsed, model.QuotaInt(5))
}
