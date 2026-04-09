package stats

import (
	"context"
	"iter"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func (repo repository) IterateStatsForDailyDecay(
	ctx context.Context,
	batchCount int32,
) iter.Seq2[persist.DailyDecayStats, error] {
	batchCount64 := int64(batchCount)

	return func(yield func(persist.DailyDecayStats, error) bool) {
		var lastIngestID model.IngestIDType = 0
		for {
			rows, err := repo.queries.ListStatsAfterIngestID(ctx, gen.ListStatsAfterIngestIDParams{
				LastIngestID: lastIngestID,
				MaxCount:     batchCount64,
			})
			if err != nil {
				yield(persist.DailyDecayStats{}, shared.MapSQLiteError("IterateStatsForDailyDecay", err))
				return
			}
			if len(rows) == 0 {
				return
			}

			for _, row := range rows {
				stats := persist.DailyDecayStats{
					IngestID:     row.IngestID,
					Score:        row.Score,
					LastViewedAt: shared.OptionTimeFromSql(row.LastViewedAt),
				}
				if !yield(stats, nil) {
					return
				}
			}
			lastIngestID = rows[len(rows)-1].IngestID
		}
	}
}
