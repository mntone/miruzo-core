package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

const (
	loveScoreThreshold = model.ScoreType(180)
	loveDayStartOffset = 5 * time.Hour
)

type testLoveAction struct {
	kind   model.ActionType
	offset time.Duration
}

func durationOptionToTime(base time.Time, d mo.Option[time.Duration]) mo.Option[time.Time] {
	if value, present := d.Get(); present {
		return mo.Some(base.Add(value))
	}
	return mo.None[time.Time]()
}

// --- ApplyLove ---

func TestStatsRepositoryApplyLoveUpdatesTimestamps(t *testing.T) {
	tests := []struct {
		name        string
		score       model.ScoreType
		loveOffset  time.Duration
		firstOffset mo.Option[time.Duration]
		lastOffset  mo.Option[time.Duration]
		wantFirst   time.Duration
		wantLast    time.Duration
	}{
		{
			name:        "Empty",
			score:       100,
			loveOffset:  20 * time.Minute,
			firstOffset: mo.None[time.Duration](),
			lastOffset:  mo.None[time.Duration](),
			wantFirst:   20 * time.Minute,
			wantLast:    20 * time.Minute,
		},
		{
			name:        "LastLovedBeforePeriodStart",
			score:       100,
			loveOffset:  35 * time.Minute,
			firstOffset: mo.Some(-48*time.Hour + 24*time.Minute),
			lastOffset:  mo.Some(-time.Microsecond),
			wantFirst:   -48*time.Hour + 24*time.Minute,
			wantLast:    35 * time.Minute,
		},
		{
			name:        "ScoreJustBelowThreshold",
			score:       loveScoreThreshold - 1,
			loveOffset:  45 * time.Minute,
			firstOffset: mo.Some(-48*time.Hour + 24*time.Minute),
			lastOffset:  mo.Some(-24*time.Hour + 24*time.Minute),
			wantFirst:   -48*time.Hour + 24*time.Minute,
			wantLast:    45 * time.Minute,
		},
	}

	ingest := mb.Ingest().Build()
	scoreDelta := model.ScoreType(20)
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
					ops.MustAddIngest(t, ingest)
					ops.MustAddStats(t, mb.
						Stats(ingest.ID).
						Score(tt.score).
						FirstLovedOffset(tt.firstOffset).
						LastLovedOffset(tt.lastOffset).
						Build())

					loveStats, err := ops.Stats().ApplyLove(
						t.Context(),
						ingest.ID,
						scoreDelta,
						baseTime.Add(tt.loveOffset),
						loveScoreThreshold,
						baseTime,
					)
					assert.NilError(t, "ApplyLove() error", err)
					assert.Equal(t, "ApplyLove().Score", loveStats.Score, tt.score+scoreDelta)
					assert.Equal(t, "ApplyLove().FirstLovedAt", loveStats.FirstLovedAt.MustGet(), baseTime.Add(tt.wantFirst))
					assert.Equal(t, "ApplyLove().LastLovedAt", loveStats.LastLovedAt.MustGet(), baseTime.Add(tt.wantLast))
				})
			})
		}
	})
}

func TestStatsRepositoryApplyLoveReturnsConflict(t *testing.T) {
	tests := []struct {
		name        string
		score       model.ScoreType
		loveOffset  time.Duration
		firstOffset mo.Option[time.Duration]
		lastOffset  mo.Option[time.Duration]
	}{
		{
			name:        "AlreadyLovedToday",
			score:       100,
			loveOffset:  24*time.Hour + 2*time.Minute,
			firstOffset: mo.Some(0 * time.Hour),
			lastOffset:  mo.Some(24 * time.Hour),
		},
		{
			name:       "ScoreAboveThreshold",
			score:      loveScoreThreshold,
			loveOffset: 24*time.Hour + 45*time.Minute,
		},
	}

	ingest := mb.Ingest().Build()
	baseTime := mb.GetDefaultBaseTime()
	scoreDelta := model.ScoreType(20)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
					ops.MustAddIngest(t, ingest)
					ops.MustAddStats(t, mb.
						Stats(ingest.ID).
						Score(tt.score).
						FirstLovedOffset(tt.firstOffset).
						LastLovedOffset(tt.lastOffset).
						Build())

					_, err := ops.Stats().ApplyLove(
						t.Context(),
						ingest.ID,
						scoreDelta,
						baseTime.Add(tt.loveOffset),
						loveScoreThreshold,
						baseTime.Add(24*time.Hour),
					)
					assert.ErrorIs(t, "ApplyLove() error", err, persist.ErrConflict)
				})
			})
		}
	})
}

func TestStatsRepositoryApplyLoveReturnsConflictWithoutStats(t *testing.T) {
	ingest := mb.Ingest().Build()
	scoreDelta := model.ScoreType(20)

	baseTime := mb.GetDefaultBaseTime()
	periodStart := baseTime.Add(-24 * time.Hour)
	loveAt := baseTime.Add(20 * time.Minute)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			_, err := ops.Stats().ApplyLove(
				t.Context(),
				ingest.ID,
				scoreDelta,
				loveAt,
				loveScoreThreshold,
				periodStart,
			)
			assert.ErrorIs(t, "ApplyLove() error", err, persist.ErrConflict)
		})
	})
}

// --- ApplyLoveCanceled ---

func TestStatsRepositoryApplyLoveCanceledUpdatesTimestamps(t *testing.T) {
	tests := []struct {
		name    string
		actions []testLoveAction
		// periodDayOffset sets periodStart as a day offset from baseTime.
		periodDayOffset  int32
		loveCancelOffset time.Duration
		firstOffset      time.Duration
		lastOffset       time.Duration
		wantFirst        mo.Option[time.Duration]
		wantLast         mo.Option[time.Duration]
	}{
		{
			name: "SingleLove",
			actions: []testLoveAction{
				// -- Day 1 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 20 * time.Minute,
				},
			},
			periodDayOffset:  0,
			loveCancelOffset: 40 * time.Minute,
			firstOffset:      20 * time.Minute,
			lastOffset:       20 * time.Minute,
			wantFirst:        mo.None[time.Duration](),
			wantLast:         mo.None[time.Duration](),
		},
		{
			name: "SingleLoveJustBeforeCancelBoundary",
			actions: []testLoveAction{
				// -- Day 1 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 20 * time.Minute,
				},
			},
			periodDayOffset:  0,
			loveCancelOffset: 20*time.Minute + time.Microsecond,
			firstOffset:      20 * time.Minute,
			lastOffset:       20 * time.Minute,
			wantFirst:        mo.None[time.Duration](),
			wantLast:         mo.None[time.Duration](),
		},
		{
			name: "EachPeriodKeepsFirst",
			actions: []testLoveAction{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 40 * time.Minute,
				},
				// -- Day 2 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 35*time.Minute,
				},
			},
			periodDayOffset:  1,
			loveCancelOffset: 24*time.Hour + 55*time.Minute,
			firstOffset:      40 * time.Minute,
			lastOffset:       24*time.Hour + 35*time.Minute,
			wantFirst:        mo.Some(40 * time.Minute),
			wantLast:         mo.Some(40 * time.Minute),
		},
		{
			name: "IncludesCancelAction",
			actions: []testLoveAction{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 10 * time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 60 * time.Minute,
				},
				// -- Day 2 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 75*time.Minute,
				},
			},
			periodDayOffset:  1,
			loveCancelOffset: 24*time.Hour + 95*time.Minute,
			firstOffset:      24*time.Hour + 75*time.Minute,
			lastOffset:       24*time.Hour + 75*time.Minute,
			wantFirst:        mo.None[time.Duration](),
			wantLast:         mo.None[time.Duration](),
		},
		{
			name: "CrossPeriodWithCancel",
			actions: []testLoveAction{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 15 * time.Minute,
				},
				// -- Day 2
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 45*time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 24*time.Hour + 50*time.Minute,
				},
				// -- Day 3 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 2*24*time.Hour + 35*time.Minute,
				},
			},
			periodDayOffset:  2,
			loveCancelOffset: 2*24*time.Hour + 50*time.Minute,
			firstOffset:      15 * time.Minute,
			lastOffset:       2*24*time.Hour + 35*time.Minute,
			wantFirst:        mo.Some(15 * time.Minute),
			wantLast:         mo.Some(15 * time.Minute),
		},
		{
			name: "DayBoundaryCancelAtBoundaryKeepsPreviousLove",
			actions: []testLoveAction{
				// -- Day 1 (previous period candidate)
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour - time.Microsecond,
				},
				// Day boundary cancel: should not cancel the previous-period candidate.
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 24 * time.Hour,
				},
				// -- Day 2 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 20*time.Minute,
				},
			},
			periodDayOffset:  1,
			loveCancelOffset: 24*time.Hour + 30*time.Minute,
			firstOffset:      24*time.Hour - time.Microsecond,
			lastOffset:       24*time.Hour + 20*time.Minute,
			wantFirst:        mo.Some(24*time.Hour - time.Microsecond),
			wantLast:         mo.Some(24*time.Hour - time.Microsecond),
		},
		{
			name: "CrossPeriodKeepsMiddleLast",
			actions: []testLoveAction{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 72 * time.Minute,
				},
				// -- Day 2
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 36*time.Minute,
				},
				// -- Day 3 (current period)
				{
					kind:   model.ActionTypeLove,
					offset: 2*24*time.Hour + 18*time.Minute,
				},
			},
			periodDayOffset:  2,
			loveCancelOffset: 2*24*time.Hour + 36*time.Minute,
			firstOffset:      72 * time.Minute,
			lastOffset:       2*24*time.Hour + 18*time.Minute,
			wantFirst:        mo.Some(72 * time.Minute),
			wantLast:         mo.Some(24*time.Hour + 36*time.Minute),
		},
	}

	ingest := mb.Ingest().Build()
	scoreDelta := model.ScoreType(-18)
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
					ops.MustAddIngest(t, ingest)
					stats := ops.MustAddStats(t, mb.
						Stats(ingest.ID).
						Score(120).
						FirstLovedOffset(tt.firstOffset).
						LastLovedOffset(tt.lastOffset).
						Build())

					for _, action := range tt.actions {
						ops.MustAddAction(
							t,
							ingest.ID,
							action.kind,
							baseTime.Add(action.offset),
						)
					}

					loveStats, err := ops.Stats().ApplyLoveCanceled(
						t.Context(),
						ingest.ID,
						scoreDelta,
						baseTime.Add(tt.loveCancelOffset),
						baseTime.Add(time.Duration(tt.periodDayOffset)*24*time.Hour),
						loveDayStartOffset,
					)
					assert.NilError(t, "ApplyLoveCanceled() error", err)
					assert.Equal(t, "ApplyLoveCanceled().Score", loveStats.Score, stats.Score+scoreDelta)
					assert.Equal(
						t,
						"ApplyLoveCanceled().FirstLovedAt",
						loveStats.FirstLovedAt,
						durationOptionToTime(baseTime, tt.wantFirst),
					)
					assert.Equal(
						t,
						"ApplyLoveCanceled().LastLovedAt",
						loveStats.LastLovedAt,
						durationOptionToTime(baseTime, tt.wantLast),
					)
				})
			})
		}
	})
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflict(t *testing.T) {
	tests := []struct {
		name    string
		actions []testLoveAction
		// periodDayOffset sets periodStart as a day offset from baseTime.
		periodDayOffset  int32
		loveCancelOffset time.Duration
		firstOffset      mo.Option[time.Duration]
		lastOffset       mo.Option[time.Duration]
	}{
		{
			name:            "NoActions",
			actions:         []testLoveAction{},
			periodDayOffset: 1,
		},
		{
			name:             "StatsHasLastLoveButNoActionRow",
			actions:          []testLoveAction{},
			periodDayOffset:  0,
			loveCancelOffset: 20 * time.Minute,
			// Defensive case: stats/action inconsistency should fail safely.
			firstOffset: mo.Some(10 * time.Minute),
			lastOffset:  mo.Some(10 * time.Minute),
		},
		{
			name: "AlreadyCanceled",
			actions: []testLoveAction{
				{
					kind:   model.ActionTypeLove,
					offset: 30 * time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 60 * time.Minute,
				},
			},
			periodDayOffset:  1,
			loveCancelOffset: 90 * time.Minute,
		},
		{
			name: "DayBoundaryCancelBeforeBoundary",
			actions: []testLoveAction{
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour - time.Microsecond,
				},
			},
			periodDayOffset:  1,
			loveCancelOffset: 24 * time.Hour,
			firstOffset:      mo.Some(24*time.Hour - time.Microsecond),
			lastOffset:       mo.Some(24*time.Hour - time.Microsecond),
		},
		{
			name:             "EqualTimestampReturnsConflict",
			actions:          []testLoveAction{},
			periodDayOffset:  0,
			loveCancelOffset: 20 * time.Minute,
			firstOffset:      mo.Some(20 * time.Minute),
			lastOffset:       mo.Some(20 * time.Minute),
		},
		{
			name: "InvalidWindowPeriodStartAfterCancel",
			actions: []testLoveAction{
				{
					kind:   model.ActionTypeLove,
					offset: 10 * time.Minute,
				},
			},
			// period_start_at (base + 24h) is after love_cancel_at (base + 20m).
			periodDayOffset:  1,
			loveCancelOffset: 20 * time.Minute,
			firstOffset:      mo.Some(10 * time.Minute),
			lastOffset:       mo.Some(10 * time.Minute),
		},
		{
			name: "IgnorePreviousPeriod",
			actions: []testLoveAction{
				{
					kind:   model.ActionTypeLove,
					offset: 30 * time.Minute,
				},
			},
			periodDayOffset:  2,
			loveCancelOffset: 24*time.Hour + 10*time.Minute,
			firstOffset:      mo.Some(30 * time.Minute),
			lastOffset:       mo.Some(30 * time.Minute),
		},
		{
			name: "IgnoresPreviousPeriodCandidates",
			actions: []testLoveAction{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 6 * time.Minute,
				},
				// -- Day 2
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 12*time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 24*time.Hour + 24*time.Minute,
				},
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 48*time.Minute,
				},
				// -- Day 3 (current period)
			},
			periodDayOffset:  3,
			loveCancelOffset: 2*24*time.Hour + 6*time.Minute,
			firstOffset:      mo.Some(6 * time.Minute),
			lastOffset:       mo.Some(24*time.Hour + 48*time.Minute),
		},
	}

	ingest := mb.Ingest().Build()
	scoreDelta := model.ScoreType(-18)
	baseTime := mb.GetDefaultBaseTime()

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
					ops.MustAddIngest(t, ingest)
					ops.MustAddStats(t, mb.
						Stats(ingest.ID).
						FirstLovedOffset(tt.firstOffset).
						LastLovedOffset(tt.lastOffset).
						Build())

					for _, action := range tt.actions {
						ops.MustAddAction(
							t,
							ingest.ID,
							action.kind,
							baseTime.Add(action.offset),
						)
					}

					_, err := ops.Stats().ApplyLoveCanceled(
						t.Context(),
						ingest.ID,
						scoreDelta,
						baseTime.Add(tt.loveCancelOffset),
						baseTime.Add(time.Duration(tt.periodDayOffset)*24*time.Hour),
						loveDayStartOffset,
					)
					assert.ErrorIs(t, "ApplyLoveCanceled() error", err, persist.ErrConflict)
				})
			})
		}
	})
}

func TestStatsRepositoryApplyLoveCanceledReturnsConflictWithoutStats(t *testing.T) {
	ingest := mb.Ingest().Build()
	scoreDelta := model.ScoreType(-18)

	baseTime := mb.GetDefaultBaseTime()
	periodStart := baseTime.Add(-24 * time.Hour)
	loveCancelAt := baseTime.Add(20 * time.Minute)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			_, err := ops.Stats().ApplyLoveCanceled(
				t.Context(),
				ingest.ID,
				scoreDelta,
				loveCancelAt,
				periodStart,
				loveDayStartOffset,
			)
			assert.ErrorIs(t, "ApplyLoveCanceled() error", err, persist.ErrConflict)
		})
	})
}
