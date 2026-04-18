package stats

import (
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

func MapStats(r gen.Stat) model.Stats {
	return model.Stats{
		IngestID:         r.IngestID,
		Score:            r.Score,
		ScoreEvaluated:   r.ScoreEvaluated,
		ScoreEvaluatedAt: persistshared.OptionTimeFromSql(r.ScoreEvaluatedAt),

		LastViewedAt: persistshared.OptionTimeFromSql(r.LastViewedAt),
		FirstLovedAt: persistshared.OptionTimeFromSql(r.FirstLovedAt),
		LastLovedAt:  persistshared.OptionTimeFromSql(r.LastLovedAt),
		HallOfFameAt: persistshared.OptionTimeFromSql(r.HallOfFameAt),

		ViewCount:               r.ViewCount,
		ViewMilestoneCount:      r.ViewMilestoneCount,
		ViewMilestoneArchivedAt: persistshared.OptionTimeFromSql(r.ViewMilestoneArchivedAt),
	}
}
