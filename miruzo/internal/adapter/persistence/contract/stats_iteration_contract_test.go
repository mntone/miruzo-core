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

func TestStatsListRepositoryIterateStatsPaginatesWithoutDuplicates(t *testing.T) {
	statsBuilders := []func(ingestID model.IngestIDType) model.Stats{
		func(ingestID model.IngestIDType) model.Stats {
			return mb.Stats(ingestID).Score(100).Build()
		},
		func(ingestID model.IngestIDType) model.Stats {
			return mb.Stats(ingestID).Score(120).ViewedOffset(2, 2*time.Hour).Build()
		},
		func(ingestID model.IngestIDType) model.Stats {
			return mb.Stats(ingestID).Score(88).ViewedOffset(4, 4*time.Hour).Build()
		},
		func(ingestID model.IngestIDType) model.Stats {
			return mb.Stats(ingestID).Build()
		},
	}

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingestIDs := make([]model.IngestIDType, len(statsBuilders))
			statsEntries := make([]model.Stats, len(statsBuilders))
			for i, statsBuilder := range statsBuilders {
				ingests := ops.MustAddIngest(t, mb.Ingest().Build())
				ingestIDs[i] = ingests.ID
				statsEntries[i] = ops.MustAddStats(t, statsBuilder(ingests.ID))
			}

			i := 0
			gotIngestIDs := make([]model.IngestIDType, len(statsEntries))
			gotRows := make([]persist.DailyDecayStats, len(statsEntries))
			for row, err := range ops.Stats().IterateStatsForDailyDecay(t.Context(), 1) {
				assert.NilError(t, "IterateStatsForDailyDecay() error", err)

				gotIngestIDs[i] = row.IngestID
				gotRows[i] = row
				i += 1
			}

			assert.EqualSlice(t, "IterateStatsForDailyDecay() ingest IDs", gotIngestIDs, ingestIDs)

			assert.Equal(t, "rows[0].Score", gotRows[0].Score, statsEntries[0].Score)
			assert.IsAbsent(t, "rows[0].LastViewedAt", gotRows[0].LastViewedAt)

			assert.Equal(t, "rows[1].Score", gotRows[1].Score, statsEntries[1].Score)
			assert.EqualFn(t, "rows[1].LastViewedAt", gotRows[1].LastViewedAt, statsEntries[1].LastViewedAt)

			assert.Equal(t, "rows[2].Score", gotRows[2].Score, statsEntries[2].Score)
			assert.EqualFn(t, "rows[2].LastViewedAt", gotRows[2].LastViewedAt, statsEntries[2].LastViewedAt)

			assert.Equal(t, "rows[3].Score", gotRows[3].Score, statsEntries[3].Score)
			assert.IsAbsent(t, "rows[3].LastViewedAt", gotRows[3].LastViewedAt)
		})
	})
}

func TestStatsListRepositoryIterateStatsReturnsEmpty(t *testing.T) {
	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			i := 0
			for range ops.Stats().IterateStatsForDailyDecay(t.Context(), 1) {
				i += 1
			}
			assert.Equal(t, "IterateStatsForDailyDecay() row count", i, 0)
		})
	})
}
