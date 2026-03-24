package stats

import (
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

func MapStats(r gen.Stat) model.Stats {
	return model.Stats{
		IngestID:         r.IngestID,
		Score:            r.Score,
		ScoreEvaluated:   r.ScoreEvaluated,
		ScoreEvaluatedAt: mo.PointerToOption(r.ScoreEvaluatedAt),

		LastViewedAt: mo.PointerToOption(r.LastViewedAt),
		FirstLovedAt: mo.PointerToOption(r.FirstLovedAt),
		LastLovedAt:  mo.PointerToOption(r.LastLovedAt),
		HallOfFameAt: mo.PointerToOption(r.HallOfFameAt),

		ViewCount:               r.ViewCount,
		ViewMilestoneCount:      r.ViewMilestoneCount,
		ViewMilestoneArchivedAt: mo.PointerToOption(r.ViewMilestoneArchivedAt),
	}
}
