package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

func TestStatsRepositoryApplyViewUpdates(t *testing.T) {
	ingest := mb.Ingest().Build()
	evaluatedAt := mb.GetDefaultBaseTime().Add(20 * time.Minute)
	scoreDelta := model.ScoreType(5)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			stats := ops.MustAddStats(t, mb.Stats(ingest.ID).Build())

			err := ops.Stats().ApplyView(t.Context(), ingest.ID, scoreDelta, evaluatedAt)
			assert.NilError(t, "ApplyView() error", err)

			dbstats := ops.MustGetStats(t, ingest.ID)
			assert.Equal(t, "stats.Score", dbstats.Score, stats.Score+scoreDelta)
			assert.IsPresent(t, "stats.LastViewedAt", dbstats.LastViewedAt)
			assert.Equal(t, "stats.LastViewedAt", dbstats.LastViewedAt.MustGet(), evaluatedAt)
			assert.Equal(t, "stats.ViewCount", dbstats.ViewCount, stats.ViewCount+1)
		})
	})
}

func TestStatsRepositoryApplyViewReturnsNotFound(t *testing.T) {
	ingest := mb.Ingest().Build()
	evaluatedAt := mb.GetDefaultBaseTime().Add(20 * time.Minute)
	scoreDelta := model.ScoreType(5)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			err := ops.Stats().ApplyView(t.Context(), ingest.ID, scoreDelta, evaluatedAt)
			assert.ErrorIs(t, "ApplyView() error", err, persist.ErrNotFound)
		})
	})
}

func TestStatsRepositoryApplyViewWithMilestoneUpdates(t *testing.T) {
	ingest := mb.Ingest().Build()
	evaluatedAt := mb.GetDefaultBaseTime().Add(25 * time.Minute)
	scoreDelta := model.ScoreType(7)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			stats := ops.MustAddStats(t, mb.
				Stats(ingest.ID).
				ViewedOffset(23, -15*time.Minute).
				Build())

			err := ops.Stats().ApplyViewWithMilestone(t.Context(), ingest.ID, scoreDelta, evaluatedAt)
			assert.NilError(t, "ApplyViewWithMilestone() error", err)

			dbstats := ops.MustGetStats(t, ingest.ID)
			assert.Equal(t, "stats.Score", dbstats.Score, stats.Score+scoreDelta)
			assert.IsPresent(t, "stats.LastViewedAt", dbstats.LastViewedAt)
			assert.Equal(t, "stats.LastViewedAt", dbstats.LastViewedAt.MustGet(), evaluatedAt)
			assert.Equal(t, "stats.ViewCount", dbstats.ViewCount, stats.ViewCount+1)
			assert.Equal(t, "stats.ViewMilestoneCount", dbstats.ViewMilestoneCount, stats.ViewCount+1)
			assert.IsPresent(t, "stats.ViewMilestoneArchivedAt", dbstats.ViewMilestoneArchivedAt)
			assert.Equal(
				t,
				"stats.ViewMilestoneArchivedAt",
				dbstats.ViewMilestoneArchivedAt.MustGet(),
				evaluatedAt,
			)
		})
	})
}

func TestStatsRepositoryApplyViewWithMilestoneReturnsNotFound(t *testing.T) {
	ingest := mb.Ingest().Build()
	evaluatedAt := mb.GetDefaultBaseTime().Add(25 * time.Minute)
	scoreDelta := model.ScoreType(7)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ops.MustAddIngest(t, ingest)
			err := ops.Stats().ApplyViewWithMilestone(t.Context(), ingest.ID, scoreDelta, evaluatedAt)
			assert.ErrorIs(t, "ApplyViewWithMilestone() error", err, persist.ErrNotFound)
		})
	})
}
