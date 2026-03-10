package stats

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func MapStats(r gen.Stat) persist.Stats {
	return persist.Stats{
		IngestID:         r.IngestID,
		Score:            model.ScoreType(r.Score),
		ScoreEvaluated:   model.ScoreType(r.ScoreEvaluated),
		ScoreEvaluatedAt: shared.OptionTimeFromSql(r.ScoreEvaluatedAt),

		LastViewedAt: shared.OptionTimeFromSql(r.LastViewedAt),
		FirstLovedAt: shared.OptionTimeFromSql(r.FirstLovedAt),
		LastLovedAt:  shared.OptionTimeFromSql(r.LastLovedAt),
		HallOfFameAt: shared.OptionTimeFromSql(r.HallOfFameAt),

		ViewCount:               r.ViewCount,
		ViewMilestoneCount:      r.ViewMilestoneCount,
		ViewMilestoneArchivedAt: shared.OptionTimeFromSql(r.ViewMilestoneArchivedAt),
	}
}
