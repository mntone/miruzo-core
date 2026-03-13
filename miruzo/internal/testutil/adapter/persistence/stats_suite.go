package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

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
	baseStats := ste.Operations.MustAddStat(t, NewStatFixture(ingest.ID))
	scoreDelta := model.ScoreType(20)
	lovedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)

	loveStats, err := ste.Repository.ApplyLove(
		ste.Context,
		ingest.ID,
		scoreDelta,
		lovedAt,
		statsSuitePeriodStartTimeUTC,
	)
	assert.NilError(t, "ApplyLove() error", err)
	assert.Equal(t, "ApplyLove().Score", loveStats.Score, baseStats.Score+scoreDelta)
	assert.Equal(t, "ApplyLove().FirstLovedAt", loveStats.FirstLovedAt.MustGet(), lovedAt)
	assert.Equal(t, "ApplyLove().LastLovedAt", loveStats.LastLovedAt.MustGet(), lovedAt)
}

func (ste StatsSuite) RunTestApplyLoveRejectsCurrentPeriod(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	ste.Operations.MustAddStat(t, NewStatFixtureWithLastLovedAt(ingest.ID, statsSuitePeriodStartTimeUTC.Add(2*time.Hour)))
	scoreDelta := model.ScoreType(20)
	lovedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)

	_, err := ste.Repository.ApplyLove(
		ste.Context,
		ingest.ID,
		scoreDelta,
		lovedAt,
		statsSuitePeriodStartTimeUTC,
	)
	assert.ErrorIs(t, "ApplyLove() error", err, persist.ErrConflict)
}

// --- view ---

func (ste StatsSuite) RunTestApplyView(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	baseStats := ste.Operations.MustAddStat(t, NewStatFixture(ingest.ID))
	evaluatedAt := statsSuiteBaseTimeUTC.Add(20 * time.Minute)
	scoreDelta := model.ScoreType(5)

	err := ste.Repository.ApplyView(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.NilError(t, "ApplyView() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStats(ste.Context, ingest.ID)
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

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, statsSuiteBaseTimeUTC))
	baseStats := ste.Operations.MustAddStat(t, NewStatFixtureWithLastViewedAt(
		ingest.ID,
		23,
		statsSuiteBaseTimeUTC.Add(-15*time.Minute),
	))
	evaluatedAt := statsSuiteBaseTimeUTC.Add(25 * time.Minute)
	scoreDelta := model.ScoreType(7)

	err := ste.Repository.ApplyViewWithMilestone(ste.Context, ingest.ID, scoreDelta, evaluatedAt)
	assert.NilError(t, "ApplyViewWithMilestone() error", err)

	imageWithStats, err := ste.ViewRepository.GetImageWithStats(ste.Context, ingest.ID)
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
