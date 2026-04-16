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
) (service *job.DailyDecayService, prov *stub.PersistenceProvider, resolver period.DailyResolver) {
	prov = stub.NewStubPersistenceProvider(dailyLoveUsed, statsEntries...)
	resolver = period.NewDailyResolver(5 * time.Hour)
	service = job.NewDailyDecay(
		prov,
		clock.NewFixedProvider(evaluatedAt),
		resolver,
		domain.NewTestScoreCalculator(resolver),
	)
	return
}

func TestDailyDecayServiceApplyUpdatesScoresAndResetsDailyLoveUsed(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	evaluatedAt := baseTime.Add(1 * time.Second)
	service, prov, resolv := newDailyDecayServiceFixture(
		3,
		evaluatedAt,

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
	_, _ = prov.ActionStub.Create(context.Background(), 4, model.ActionTypeDecay, evaluatedAt, resolv.PeriodStart(baseTime))

	err := service.ApplyDailyDecay(context.Background())
	assert.NilError(t, "ApplyDailyDecay() error", err)
	assert.Equal(t, "IterateStatsForDailyDecay() batchCount", prov.StatsStub.IterateStatsForDailyDecayArgs[0], 500)

	assert.Equal(t, "store.Stats(IngestID=1).Score", prov.StatsStub.Store[1].Score, 80)
	assert.IsAbsent(t, "store.Stats(IngestID=1).ScoreEvaluatedAt", prov.StatsStub.Store[1].ScoreEvaluatedAt)

	assert.Equal(t, "store.Stats(IngestID=2).Score", prov.StatsStub.Store[2].Score, 120-2)
	assert.EqualFn(t, "store.Stats(IngestID=2).ScoreEvaluatedAt", prov.StatsStub.Store[2].ScoreEvaluatedAt, mo.Some(evaluatedAt))

	assert.Equal(t, "store.Stats(IngestID=3).Score", prov.StatsStub.Store[3].Score, 180-3)
	assert.EqualFn(t, "store.Stats(IngestID=3).ScoreEvaluatedAt", prov.StatsStub.Store[3].ScoreEvaluatedAt, mo.Some(evaluatedAt))

	assert.Equal(t, "store.Stats(IngestID=4).Score", prov.StatsStub.Store[4].Score, 138)
	assert.EqualFn(t, "store.Stats(IngestID=4).ScoreEvaluatedAt", prov.StatsStub.Store[4].ScoreEvaluatedAt, mo.Some(baseTime))

	assert.LenIs(t, "store.Action count", prov.ActionStub.Store, 3)
	assert.Equal(t, "store.User.DailyLoveUsed", prov.UserStub.DailyLoveUsed, 0)
}

func TestDailyDecayServiceApplyReturnsIterateError(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	service, mgr, _ := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	expectedErr := errors.New("iterate failed")
	mgr.StatsStub.IterateStatsForDailyDecayError = expectedErr

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, expectedErr)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.UserStub.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplySkipsOnCreateDailyDecayIfAbsentConflict(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	service, mgr, _ := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.ActionStub.CreateDailyDecayIfAbsentError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.NilError(t, "ApplyDailyDecay() error", err)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.StatsStub.Store[1].Score, model.ScoreType(120))
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.UserStub.DailyLoveUsed, model.QuotaInt(0))
}

func TestDailyDecayServiceApplyReturnsCreateDailyDecayIfAbsentError(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	service, mgr, _ := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.ActionStub.CreateDailyDecayIfAbsentError = persist.ErrConnectionLost

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrServiceUnavailable)
	assert.Empty(t, "store.Action", mgr.ActionStub.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.UserStub.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsApplyDecayError(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	service, mgr, _ := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.StatsStub.ApplyDecayError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.StatsStub.Store[1].Score, model.ScoreType(120))
	assert.Empty(t, "store.Action", mgr.ActionStub.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.UserStub.DailyLoveUsed, model.QuotaInt(5))
}

func TestDailyDecayServiceApplyReturnsResetDailyLoveUsedError(t *testing.T) {
	baseTime := mb.GetDefaultBaseTime()
	service, mgr, _ := newDailyDecayServiceFixture(
		5,
		baseTime.Add(1*time.Second),
		mb.Stats(1).Score(120).Viewed(2, baseTime.Add(-time.Hour)).Build(),
	)
	mgr.UserStub.ResetError = persist.ErrConflict

	err := service.ApplyDailyDecay(context.Background())
	assert.ErrorIs(t, "ApplyDailyDecay() error", err, serviceerror.ErrConflict)
	assert.Equal(t, "store.Stats(IngestID=1).Score", mgr.StatsStub.Store[1].Score, model.ScoreType(120))
	assert.Empty(t, "store.Action", mgr.ActionStub.Store)
	assert.Equal(t, "store.User.DailyLoveUsed", mgr.UserStub.DailyLoveUsed, model.QuotaInt(5))
}
