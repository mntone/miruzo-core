package persistence

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

func (ste StatsSuite) RunTestIterateStatsPaginatesWithoutDuplicates(t *testing.T) {
	t.Helper()

	baseTime := mb.GetDefaultStatsBaseTime()
	statsEntries := []model.Stats{
		mb.Stats(1).Score(100).Build(),
		mb.Stats(2).Score(120).Viewed(2, baseTime.Add(2*time.Hour)).Build(),
		mb.Stats(4).Score(88).Viewed(4, baseTime.Add(4*time.Hour)).Build(),
		mb.Stats(5).Build(),
	}

	ingestIDs := make([]model.IngestIDType, len(statsEntries))
	for i, stats := range statsEntries {
		ingestIDs[i] = stats.IngestID
		ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(stats.IngestID, baseTime))
		ste.Operations.MustAddStat(t, stats)
	}

	i := 0
	gotIngestIDs := make([]model.IngestIDType, len(statsEntries))
	gotRows := make(map[model.IngestIDType]persist.DailyDecayStats, len(statsEntries))
	for row, err := range ste.Repository.IterateStatsForDailyDecay(ste.Context, 1) {
		assert.NilError(t, "IterateStatsForDailyDecay() error", err)

		gotIngestIDs[i] = row.IngestID
		gotRows[row.IngestID] = row
		i += 1
	}

	assert.EqualSlice(t, "IterateStatsForDailyDecay() ingest IDs", gotIngestIDs, ingestIDs)

	assert.Equal(t, "rows(IngestID=1).Score", gotRows[1].Score, statsEntries[0].Score)
	assert.IsAbsent(t, "rows(IngestID=1).LastViewedAt", gotRows[1].LastViewedAt)

	assert.Equal(t, "rows(IngestID=2).Score", gotRows[2].Score, statsEntries[1].Score)
	assert.EqualFn(t, "rows(IngestID=2).LastViewedAt", gotRows[2].LastViewedAt, statsEntries[1].LastViewedAt)

	assert.Equal(t, "rows(IngestID=4).Score", gotRows[4].Score, statsEntries[2].Score)
	assert.EqualFn(t, "rows(IngestID=4).LastViewedAt", gotRows[4].LastViewedAt, statsEntries[2].LastViewedAt)

	assert.Equal(t, "rows(IngestID=5).Score", gotRows[5].Score, statsEntries[3].Score)
	assert.IsAbsent(t, "rows(IngestID=5).LastViewedAt", gotRows[5].LastViewedAt)
}

func (ste StatsSuite) RunTestIterateStatsReturnsEmpty(t *testing.T) {
	t.Helper()

	i := 0
	for range ste.Repository.IterateStatsForDailyDecay(ste.Context, 1) {
		i += 1
	}
	assert.Equal(t, "IterateStatsForDailyDecay() row count", i, 0)
}
