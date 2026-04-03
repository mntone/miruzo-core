package persistence

import (
	"context"
	"fmt"
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

// --- schema ---

func (ste StatsSuite) RunTestStatsSchemaRejectsInvalidScore(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		score   int32
		wantErr error
	}{
		{
			name:    "score=-32769",
			score:   -32769,
			wantErr: persist.ErrCheckViolation,
		},
		{
			name:    "score=32768",
			score:   32768,
			wantErr: persist.ErrCheckViolation,
		},
	}

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	for _, tt := range tests {
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO stats(ingest_id, score, score_evaluated) VALUES(%d, %d, 100)",
				ingest.ID, tt.score,
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}

func (ste StatsSuite) RunTestStatsSchemaRejectsInvalidScoreEvaluated(t *testing.T) {
	t.Helper()

	tests := []struct {
		name    string
		score   int32
		scoreAt time.Time
		wantErr error
	}{
		{
			name:    "score_evaluated=-32769",
			score:   -32769,
			scoreAt: statsSuiteBaseTimeUTC.Add(20 * time.Minute),
			wantErr: persist.ErrCheckViolation,
		},
		{
			name:    "score_evaluated=32768",
			score:   32768,
			scoreAt: statsSuiteBaseTimeUTC.Add(40 * time.Minute),
			wantErr: persist.ErrCheckViolation,
		},
	}

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	for _, tt := range tests {
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			stmt := fmt.Sprintf(
				"INSERT INTO stats(ingest_id, score, score_evaluated, score_evaluated_at) VALUES(%d, 100, %d, '%s')",
				ingest.ID,
				tt.score,
				tt.scoreAt.Format(time.RFC3339Nano),
			)
			err := ste.Operations.ExecuteStatement(stmt)
			assert.ErrorIs(t, "insert error", err, tt.wantErr)
		})
	}
}

// PostgreSQL only
func (ste StatsSuite) RunTestStatsSchemaRejectsInvalidOccurredAt(t *testing.T) {
	t.Helper()

	tests := []struct {
		name string
		stmt string
	}{
		{
			name: "score_evaluated_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, score_evaluated_at) VALUES(%d, 100, 100, 'infinity')",
		},
		{
			name: "first_loved_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, first_loved_at, last_loved_at) VALUES(%d, 100, 100, 'infinity', 'infinity')",
		},
		{
			name: "last_loved_at=-infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, first_loved_at, last_loved_at) VALUES(%d, 100, 100, '-infinity', '-infinity')",
		},
		{
			name: "hall_of_fame_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, hall_of_fame_at) VALUES(%d, 100, 100, 'infinity')",
		},
		{
			name: "last_viewed_at=infinity",
			stmt: "INSERT INTO stats(ingest_id, score, score_evaluated, last_viewed_at, view_count) VALUES(%d, 100, 100, 'infinity', 1)",
		},
	}

	ingest := ste.Operations.MustAddIngest(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	for _, tt := range tests {
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			err := ste.Operations.ExecuteStatement(fmt.Sprintf(tt.stmt, ingest.ID))
			assert.ErrorIs(t, "insert error", err, persist.ErrCheckViolation)
		})
	}
}

// --- daily decay ---

func (ste StatsSuite) RunTestApplyDecayUpdates(t *testing.T) {
	t.Helper()

	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, baseTime))
	stats := ste.Operations.MustAddStat(t, mb.Stats(ingest.ID).Build())

	evaluatedAt := baseTime.Add(24 * time.Hour)
	err := ste.Repository.ApplyDecay(ste.Context, ingest.ID, stats.Score-2, evaluatedAt)
	assert.NilError(t, "ApplyDecay() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStatsForUpdate() error", err)
	assert.Equal(t, "imageWithStats.Stats.Score", imageWithStats.Stats.Score, stats.Score-2)
	assert.Equal(t, "imageWithStats.Stats.ScoreEvaluated", imageWithStats.Stats.ScoreEvaluated, stats.Score-2)
	assert.IsPresent(t, "imageWithStats.Stats.ScoreEvaluatedAt", imageWithStats.Stats.ScoreEvaluatedAt)
	assert.Equal(t, "imageWithStats.Stats.ScoreEvaluatedAt", imageWithStats.Stats.ScoreEvaluatedAt.MustGet(), evaluatedAt)
}

func (ste StatsSuite) RunTestApplyDecayReturnsConflictWithoutStats(t *testing.T) {
	t.Helper()

	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, baseTime))

	evaluatedAt := baseTime.Add(24 * time.Hour)
	err := ste.Repository.ApplyDecay(
		ste.Context,
		ingest.ID,
		98,
		evaluatedAt,
	)
	assert.ErrorIs(t, "ApplyDecay() error", err, persist.ErrConflict)
}

// --- hall of fame granted ---

func (ste StatsSuite) RunTestApplyHallOfFameGrantedUpdates(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	ste.Operations.MustAddStat(t, mb.Stats(ingest.ID).Score(180).Build())
	hallOfFameAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreThreshold := model.ScoreType(180)

	err := ste.Repository.ApplyHallOfFameGranted(
		ste.Context,
		ingest.ID,
		hallOfFameAt,
		scoreThreshold,
	)
	assert.NilError(t, "ApplyHallOfFameGranted() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStatsForUpdate() error", err)
	assert.IsPresent(t, "imageWithStats.Stats.HallOfFameAt", imageWithStats.Stats.HallOfFameAt)
	assert.Equal(
		t,
		"imageWithStats.Stats.HallOfFameAt",
		imageWithStats.Stats.HallOfFameAt.MustGet(),
		hallOfFameAt,
	)
}

func (ste StatsSuite) RunTestApplyHallOfFameGrantedReturnsConflict(t *testing.T) {
	t.Helper()

	tests := []struct {
		name   string
		score  model.ScoreType
		offset mo.Option[time.Duration]
	}{
		{
			name:  "ScoreBelowThreshold",
			score: 170,
		},
		{
			name:   "AlreadyGranted",
			score:  180,
			offset: mo.Some(2 * time.Hour),
		},
	}

	// ingest
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	hallOfFameAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreThreshold := model.ScoreType(180)

	for _, tt := range tests {
		ste.Operations.MustTruncateStats(t)

		t.Run(tt.name, func(t *testing.T) {
			ste.Operations.MustAddStat(t, mb.
				Stats(ingest.ID).
				Score(tt.score).
				HallOfFameOffset(tt.offset).
				Build())

			err := ste.Repository.ApplyHallOfFameGranted(
				ste.Context,
				ingest.ID,
				hallOfFameAt,
				scoreThreshold,
			)
			assert.ErrorIs(t, "ApplyHallOfFameGranted() error", err, persist.ErrConflict)
		})
	}
}

func (ste StatsSuite) RunTestApplyHallOfFameGrantedReturnsConflictWithoutStats(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	hallOfFameAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreThreshold := model.ScoreType(180)

	err := ste.Repository.ApplyHallOfFameGranted(
		ste.Context,
		ingest.ID,
		hallOfFameAt,
		scoreThreshold,
	)
	assert.ErrorIs(t, "ApplyHallOfFameGranted() error", err, persist.ErrConflict)
}

// --- hall of fame revoked ---

func (ste StatsSuite) RunTestApplyHallOfFameRevokedUpdates(t *testing.T) {
	t.Helper()

	baseTime := mb.GetDefaultStatsBaseTime()
	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, baseTime))
	ste.Operations.MustAddStat(t, mb.
		Stats(ingest.ID).
		Score(180).
		HallOfFameOffset(20*time.Minute).
		Build())

	err := ste.Repository.ApplyHallOfFameRevoked(ste.Context, ingest.ID)
	assert.NilError(t, "ApplyHallOfFameRevoked() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStatsForUpdate() error", err)
	assert.IsAbsent(t, "imageWithStats.Stats.HallOfFameAt", imageWithStats.Stats.HallOfFameAt)
}

func (ste StatsSuite) RunTestApplyHallOfFameRevokedReturnsConflict(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	ste.Operations.MustAddStat(t, mb.
		Stats(ingest.ID).
		Score(180).
		Build())

	err := ste.Repository.ApplyHallOfFameRevoked(ste.Context, ingest.ID)
	assert.ErrorIs(t, "ApplyHallOfFameRevoked() error", err, persist.ErrConflict)
}

func (ste StatsSuite) RunTestApplyHallOfFameRevokedReturnsConflictWithoutStats(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))

	err := ste.Repository.ApplyHallOfFameRevoked(ste.Context, ingest.ID)
	assert.ErrorIs(t, "ApplyHallOfFameRevoked() error", err, persist.ErrConflict)
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

// --- view ---

func (ste StatsSuite) RunTestApplyView(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	baseStats := ste.Operations.MustAddStat(t, mb.Stats(ingest.ID).Build())
	evaluatedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreDelta := model.ScoreType(5)

	err := ste.Repository.ApplyView(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.NilError(t, "ApplyView() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStats() error", err)
	assert.Equal(t, "imageWithStats.Stats.Score", imageWithStats.Stats.Score, baseStats.Score+scoreDelta)
	assert.IsPresent(t, "imageWithStats.Stats.LastViewedAt", imageWithStats.Stats.LastViewedAt)
	assert.Equal(
		t,
		"imageWithStats.Stats.LastViewedAt",
		imageWithStats.Stats.LastViewedAt.MustGet(),
		evaluatedAt,
	)
	assert.Equal(
		t,
		"imageWithStats.Stats.ViewCount",
		imageWithStats.Stats.ViewCount,
		baseStats.ViewCount+1,
	)
}

func (ste StatsSuite) RunTestApplyViewNotFound(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	evaluatedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreDelta := model.ScoreType(5)

	err := ste.Repository.ApplyView(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.ErrorIs(t, "ApplyView() error", err, persist.ErrNotFound)
}

func (ste StatsSuite) RunTestApplyViewWithMilestone(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuitePeriodStartTimeUTC))

	baseStats := ste.Operations.MustAddStat(t, mb.
		Stats(ingest.ID).
		ChangeBaseTime(statsSuiteBaseTimeUTC).
		ViewedOffset(23, -15*time.Minute).
		Build())
	evaluatedAt := statsSuiteBaseTimeUTC.Add(25 * time.Minute)
	scoreDelta := model.ScoreType(7)

	err := ste.Repository.ApplyViewWithMilestone(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.NilError(t, "ApplyViewWithMilestone() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStats() error", err)
	assert.Equal(t, "imageWithStats.Stats.Score", imageWithStats.Stats.Score, baseStats.Score+scoreDelta)
	assert.IsPresent(t, "imageWithStats.Stats.LastViewedAt", imageWithStats.Stats.LastViewedAt)
	assert.Equal(
		t,
		"imageWithStats.Stats.LastViewedAt",
		imageWithStats.Stats.LastViewedAt.MustGet(),
		evaluatedAt,
	)
	assert.Equal(
		t,
		"imageWithStats.Stats.ViewCount",
		imageWithStats.Stats.ViewCount,
		baseStats.ViewCount+1,
	)

	assert.Equal(
		t,
		"imageWithStats.Stats.ViewMilestoneCount",
		imageWithStats.Stats.ViewMilestoneCount,
		baseStats.ViewCount+1,
	)
	assert.IsPresent(
		t,
		"imageWithStats.Stats.ViewMilestoneArchivedAt",
		imageWithStats.Stats.ViewMilestoneArchivedAt,
	)
	viewMilestoneArchivedAt, _ := imageWithStats.Stats.ViewMilestoneArchivedAt.Get()
	assert.Equal(
		t,
		"imageWithStats.Stats.ViewMilestoneArchivedAt",
		viewMilestoneArchivedAt,
		evaluatedAt,
	)
}

func (ste StatsSuite) RunTestApplyViewWithMilestoneNotFound(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	evaluatedAt := statsSuiteBaseTimeUTC.Add(25 * time.Minute)
	scoreDelta := model.ScoreType(7)

	err := ste.Repository.ApplyViewWithMilestone(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.ErrorIs(t, "ApplyViewWithMilestone() error", err, persist.ErrNotFound)
}
