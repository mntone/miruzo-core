package persistence

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

var viewSuiteBaseTimeUTC = time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

type ViewSuite SuiteBase[persist.ViewRepository]

func (ste ViewSuite) RunTestGetImageWithStats(t *testing.T) {
	t.Helper()

	ingest := ste.Operations.MustAddIngestAndImage(t, NewIngestFixture(1, viewSuiteBaseTimeUTC))
	stats := ste.Operations.MustAddStat(t, NewStatFixtureWithLastViewedAt(
		ingest.ID,
		24,
		viewSuiteBaseTimeUTC.Add(30*time.Minute),
	))

	imageWithStats, err := ste.Repository.GetImageWithStatsForUpdate(ste.Context, ingest.ID)
	assert.NilError(t, "GetImageWithStats() error", err)
	assert.Equal(t, "imageWithStats.IngestID", imageWithStats.IngestID, ingest.ID)
	assert.Equal(t, "imageWithStats.IngestedAt", imageWithStats.IngestedAt, viewSuiteBaseTimeUTC)
	assert.Equal(t, "imageWithStats.Stats.IngestID", imageWithStats.Stats.IngestID, stats.IngestID)
	assert.Equal(t, "imageWithStats.Stats.ViewCount", imageWithStats.Stats.ViewCount, stats.ViewCount)
}
