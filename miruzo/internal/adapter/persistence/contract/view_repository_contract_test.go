package contract_test

import (
	"testing"
	"time"

	c "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/contract"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	mb "github.com/mntone/miruzo-core/miruzo/internal/testutil/modelbuilder"
)

func TestViewRepositoryGetImageWithStatsForUpdate(t *testing.T) {
	baseTime := time.Date(2026, 1, 9, 15, 0, 0, 0, time.UTC)

	runHarnesses(t, func(t *testing.T, h c.Harness) {
		h.RunInTx(t, func(t *testing.T, ops c.TxSession) {
			ingest := ops.MustAddIngest(t, mb.
				Ingest().
				Ingested(baseTime).
				Captured(baseTime).
				Updated(baseTime).
				Build())
			ops.MustAddImage(t, mb.Image(ingest.ID).Ingested(baseTime).Build())
			ops.MustAddStats(t, mb.
				Stats(ingest.ID).
				ChangeBaseTime(baseTime).
				ViewedOffset(24, 30*time.Minute).
				Build())

			imageWithStats, err := ops.View().GetImageWithStatsForUpdate(t.Context(), ingest.ID)
			assert.NilError(t, "GetImageWithStatsForUpdate() error", err)
			assert.Equal(t, "imageWithStats.IngestID", imageWithStats.IngestID, ingest.ID)
			assert.Equal(t, "imageWithStats.IngestedAt", imageWithStats.IngestedAt, baseTime)
			assert.Equal(t, "imageWithStats.Stats.IngestID", imageWithStats.Stats.IngestID, ingest.ID)
			assert.Equal(t, "imageWithStats.Stats.ViewCount", imageWithStats.Stats.ViewCount, int64(24))
		})
	})
}
