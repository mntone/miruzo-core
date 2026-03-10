package stats

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func MapStats(r gen.Stat) persist.Stats {
	return persist.Stats{
		IngestID:         r.IngestID,
		Score:            r.Score,
		ScoreEvaluated:   r.ScoreEvaluated,
		ScoreEvaluatedAt: shared.OptionTimeFromPgtype(r.ScoreEvaluatedAt),

		LastViewedAt: shared.OptionTimeFromPgtype(r.LastViewedAt),
		FirstLovedAt: shared.OptionTimeFromPgtype(r.FirstLovedAt),
		LastLovedAt:  shared.OptionTimeFromPgtype(r.LastLovedAt),
		HallOfFameAt: shared.OptionTimeFromPgtype(r.HallOfFameAt),

		ViewCount:               r.ViewCount,
		ViewMilestoneCount:      r.ViewMilestoneCount,
		ViewMilestoneArchivedAt: shared.OptionTimeFromPgtype(r.ViewMilestoneArchivedAt),
	}
}
