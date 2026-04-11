package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
	"github.com/samber/mo"
)

var statsDayStartOffset = 5 * time.Hour

var statsSuiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)
var statsSuitePeriodStartTimeUTC = time.Date(2026, 1, 9, 5, 0, 0, 0, time.UTC)

type StatsSuite struct {
	Context        context.Context
	Operations     Operations
	Repository     persist.StatsRepository
	ViewRepository persist.ViewRepository
}

// --- love ---

func (ste StatsSuite) RunTestApplyLoveUpdatesWhenEmpty(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	baseStats := ste.Operations.MustAddStat(t, mb.Stats(ingest.ID).Build())
	scoreDelta := model.ScoreType(20)
	lovedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)

	loveStats, err := ste.Repository.ApplyLove(
		ste.Context,
		ingest.ID,
		scoreDelta,
		lovedAt,
		180,
		statsSuitePeriodStartTimeUTC,
	)
	assert.NilError(t, "ApplyLove() error", err)
	assert.Equal(t, "ApplyLove().Score", loveStats.Score, baseStats.Score+scoreDelta)
	assert.Equal(t, "ApplyLove().FirstLovedAt", loveStats.FirstLovedAt.MustGet(), lovedAt)
	assert.Equal(t, "ApplyLove().LastLovedAt", loveStats.LastLovedAt.MustGet(), lovedAt)
}

func (ste StatsSuite) RunTestApplyLoveReturnsConflict(t *testing.T) {
	t.Helper()

	tests := []struct {
		name   string
		period time.Duration
		score  model.ScoreType
		first  mo.Option[time.Duration]
		last   mo.Option[time.Duration]
		love   time.Duration
	}{
		{
			name:  "AlreadyLovedToday",
			score: 100,
			first: mo.Some(0 * time.Hour),
			last:  mo.Some(24 * time.Hour),
			love:  24*time.Hour + 2*time.Minute,
		},
		{
			name:  "ScoreAboveThreshold",
			score: 180,
			love:  24*time.Hour + 45*time.Minute,
		},
	}

	// ingest
	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, baseTime))

	scoreDelta := model.ScoreType(20)
	for _, tt := range tests {
		ste.Operations.MustTruncateActions(t)
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			ste.Operations.MustAddStat(t, mb.
				Stats(ingest.ID).
				Score(tt.score).
				FirstLovedOffset(tt.first).
				LastLovedOffset(tt.last).
				Build())

			_, err := ste.Repository.ApplyLove(
				ste.Context,
				ingest.ID,
				scoreDelta,
				baseTime.Add(tt.love),
				180,
				baseTime.Add(24*time.Hour),
			)
			assert.ErrorIs(t, "ApplyLove() error", err, persist.ErrConflict)
		})
	}
}

// --- love cancel ---

type testActions struct {
	kind   model.ActionType
	offset time.Duration
}

func toTime(d mo.Option[time.Duration]) mo.Option[time.Time] {
	value, present := d.Get()
	if !present {
		return mo.None[time.Time]()
	}

	return mo.Some(statsSuitePeriodStartTimeUTC.Add(value))
}

func (ste StatsSuite) RunTestApplyLoveCanceledUpdatesTimestamps(t *testing.T) {
	t.Helper()

	tests := []struct {
		name      string
		actions   []testActions
		period    time.Duration
		first     mo.Option[time.Duration]
		last      mo.Option[time.Duration]
		wantFirst mo.Option[time.Duration]
		wantLast  mo.Option[time.Duration]
	}{
		{
			name: "SingleLove",
			actions: []testActions{
				{
					kind:   model.ActionTypeLove,
					offset: 20 * time.Minute,
				},
			},
			period:    0,
			first:     mo.Some(20 * time.Minute),
			last:      mo.Some(20 * time.Minute),
			wantFirst: mo.None[time.Duration](),
			wantLast:  mo.None[time.Duration](),
		},
		{
			name: "EachPeriodKeepsFirst",
			actions: []testActions{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 40 * time.Minute,
				},
				// -- Day 2
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 35*time.Minute,
				},
			},
			period:    24 * time.Hour,
			first:     mo.Some(40 * time.Minute),
			last:      mo.Some(24*time.Hour + 35*time.Minute),
			wantFirst: mo.Some(40 * time.Minute),
			wantLast:  mo.Some(40 * time.Minute),
		},
		{
			name: "IncludesCancelAction",
			actions: []testActions{
				// -- Day 1
				{
					kind:   model.ActionTypeLove,
					offset: 10 * time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 60 * time.Minute,
				},
				// -- Day 2
				{
					kind:   model.ActionTypeLove,
					offset: 24*time.Hour + 75*time.Minute,
				},
			},
			period:    24 * time.Hour,
			first:     mo.Some(24*time.Hour + 75*time.Minute),
			last:      mo.Some(24*time.Hour + 75*time.Minute),
			wantFirst: mo.None[time.Duration](),
			wantLast:  mo.None[time.Duration](),
		},
		{
			name: "CrossPeriodWithCancel",
			actions: []testActions{
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
				// -- Day 3
				{
					kind:   model.ActionTypeLove,
					offset: 2*24*time.Hour + 35*time.Minute,
				},
			},
			period:    2 * 24 * time.Hour,
			first:     mo.Some(15 * time.Minute),
			last:      mo.Some(2*24*time.Hour + 35*time.Minute),
			wantFirst: mo.Some(15 * time.Minute),
			wantLast:  mo.Some(15 * time.Minute),
		},
		{
			name: "CrossPeriodKeepsMiddleLast",
			actions: []testActions{
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
				// -- Day 3
				{
					kind:   model.ActionTypeLove,
					offset: 2*24*time.Hour + 18*time.Minute,
				},
			},
			period:    2 * 24 * time.Hour,
			first:     mo.Some(72 * time.Minute),
			last:      mo.Some(2*24*time.Hour + 18*time.Minute),
			wantFirst: mo.Some(72 * time.Minute),
			wantLast:  mo.Some(24*time.Hour + 36*time.Minute),
		},
	}

	// ingest
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	scoreDelta := model.ScoreType(-18)
	for _, tt := range tests {
		ste.Operations.MustTruncateActions(t)
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			// stats
			baseStats := ste.Operations.MustAddStat(t, mb.
				Stats(ingest.ID).
				ChangeBaseTime(statsSuitePeriodStartTimeUTC).
				Score(120).
				FirstLovedOffset(tt.first).
				LastLovedOffset(tt.last).
				Build())

			// actions
			for _, a := range tt.actions {
				ste.Operations.MustAddAction(
					t,
					ingest.ID,
					a.kind,
					statsSuitePeriodStartTimeUTC.Add(a.offset),
				)
			}

			// exec statement
			loveStats, err := ste.Repository.ApplyLoveCanceled(
				ste.Context,
				ingest.ID,
				scoreDelta,
				statsSuitePeriodStartTimeUTC.Add(tt.period),
				statsDayStartOffset,
			)
			assert.NilError(t, "ApplyLoveCanceled() error", err)
			assert.Equal(t, "ApplyLoveCanceled().Score", loveStats.Score, baseStats.Score+scoreDelta)
			assert.Equal(t, "ApplyLoveCanceled().FirstLovedAt", loveStats.FirstLovedAt, toTime(tt.wantFirst))
			assert.Equal(t, "ApplyLoveCanceled().LastLovedAt", loveStats.LastLovedAt, toTime(tt.wantLast))
		})
	}
}

func (ste StatsSuite) RunTestApplyLoveCanceledReturnsConflict(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		actions []testActions
		period  time.Duration
		first   mo.Option[time.Duration]
		last    mo.Option[time.Duration]
		wantErr error
	}{
		{
			name:    "NoActions",
			actions: []testActions{},
			period:  0,
			wantErr: persist.ErrConflict,
		},
		{
			name: "AlreadyCanceled",
			actions: []testActions{
				{
					kind:   model.ActionTypeLove,
					offset: 30 * time.Minute,
				},
				{
					kind:   model.ActionTypeLoveCanceled,
					offset: 60 * time.Minute,
				},
			},
			period:  0,
			wantErr: persist.ErrConflict,
		},
		{
			name: "IgnorePreviousPeriod",
			actions: []testActions{
				{
					kind:   model.ActionTypeLove,
					offset: 30 * time.Minute,
				},
			},
			period:  24 * time.Hour,
			first:   mo.Some(30 * time.Minute),
			last:    mo.Some(30 * time.Minute),
			wantErr: persist.ErrConflict,
		},
		{
			name: "DayBoundaryBeforeDayStart",
			actions: []testActions{
				{
					kind:   model.ActionTypeLove,
					offset: 23*time.Hour + 30*time.Minute,
				},
			},
			period:  24 * time.Hour,
			first:   mo.Some(23*time.Hour + 30*time.Minute),
			last:    mo.Some(23*time.Hour + 30*time.Minute),
			wantErr: persist.ErrConflict,
		},
		{
			name: "IgnoresPreviousPeriodCandidates",
			actions: []testActions{
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
			},
			period:  2 * 24 * time.Hour,
			first:   mo.Some(6 * time.Minute),
			last:    mo.Some(24*time.Hour + 48*time.Minute),
			wantErr: persist.ErrConflict,
		},
	}

	// ingest
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	scoreDelta := model.ScoreType(-18)
	for _, tt := range tests {
		ste.Operations.MustTruncateActions(t)
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			// stats
			ste.Operations.MustAddStat(t, mb.
				Stats(ingest.ID).
				ChangeBaseTime(statsSuitePeriodStartTimeUTC).
				FirstLovedOffset(tt.first).
				LastLovedOffset(tt.last).
				Build())

			// actions
			for _, a := range tt.actions {
				ste.Operations.MustAddAction(
					t,
					ingest.ID,
					a.kind,
					statsSuitePeriodStartTimeUTC.Add(a.offset),
				)
			}

			// exec statement
			_, err := ste.Repository.ApplyLoveCanceled(
				ste.Context,
				ingest.ID,
				scoreDelta,
				statsSuitePeriodStartTimeUTC.Add(tt.period),
				statsDayStartOffset,
			)
			assert.ErrorIs(t, "ApplyLoveCanceled() error", err, tt.wantErr)
		})
	}
}

func (ste StatsSuite) RunTestApplyLoveCanceledReturnsConflictWithoutStats(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	scoreDelta := model.ScoreType(-18)

	_, err := ste.Repository.ApplyLoveCanceled(
		ste.Context,
		ingest.ID,
		scoreDelta,
		statsSuitePeriodStartTimeUTC,
		statsDayStartOffset,
	)
	assert.ErrorIs(t, "ApplyLoveCanceled() error", err, persist.ErrConflict)
}
