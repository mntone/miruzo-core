package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

func TestStatsRepositoryApplyDecayUpdates(t *testing.T) {
	evaluatedAt := mb.GetDefaultBaseTime().Add(24 * time.Hour)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			stats := ops.MustAddStats(t, mb.Stats(ingest.ID).Build())

			err := ops.Stats().ApplyDecay(t.Context(), ingest.ID, stats.Score-2, evaluatedAt)
			assert.NilError(t, "ApplyDecay() error", err)

			dbstats := ops.MustGetStats(t, ingest.ID)
			assert.Equal(t, "stats.Score", dbstats.Score, stats.Score-2)
			assert.Equal(t, "stats.ScoreEvaluated", dbstats.ScoreEvaluated, stats.Score-2)
			assert.IsPresent(t, "stats.ScoreEvaluatedAt", dbstats.ScoreEvaluatedAt)
			assert.Equal(t, "stats.ScoreEvaluatedAt", dbstats.ScoreEvaluatedAt.MustGet(), evaluatedAt)
		})
	})
}

func TestStatsRepositoryApplyDecayReturnsConflictWithoutStats(t *testing.T) {
	evaluatedAt := mb.GetDefaultBaseTime().Add(24 * time.Hour)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.Ingest().Build())
			err := ops.Stats().ApplyDecay(t.Context(), ingest.ID, 98, evaluatedAt)
			assert.ErrorIs(t, "ApplyDecay() error", err, persist.ErrConflict)
		})
	})
}
