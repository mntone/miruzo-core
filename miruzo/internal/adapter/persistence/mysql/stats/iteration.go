package stats

import (
	"context"
	"iter"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/dberrors"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func (repo repository) IterateStatsForDailyDecay(
	ctx context.Context,
	batchCount int32,
) iter.Seq2[persist.DailyDecayStats, error] {
	return func(yield func(persist.DailyDecayStats, error) bool) {
		var lastIngestID model.IngestIDType = 0
		for {
			rows, err := repo.queries.ListStatsAfterIngestID(ctx, gen.ListStatsAfterIngestIDParams{
				LastIngestID: lastIngestID,
				Limit:        batchCount,
			})
			if err != nil {
				yield(persist.DailyDecayStats{}, dberrors.ToPersist("IterateStatsForDailyDecay", err))
				return
			}
			if len(rows) == 0 {
				return
			}

			for _, row := range rows {
				stats := persist.DailyDecayStats{
					IngestID:     row.IngestID,
					Score:        row.Score,
					LastViewedAt: persistshared.OptionTimeFromSql(row.LastViewedAt),
				}
				if !yield(stats, nil) {
					return
				}
			}
			lastIngestID = rows[len(rows)-1].IngestID
		}
	}
}
