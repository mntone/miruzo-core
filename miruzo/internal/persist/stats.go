package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Stats struct {
	IngestID         model.IngestIDType
	Score            model.ScoreType
	ScoreEvaluated   model.ScoreType
	ScoreEvaluatedAt mo.Option[time.Time]

	LastViewedAt mo.Option[time.Time]
	FirstLovedAt mo.Option[time.Time]
	LastLovedAt  mo.Option[time.Time]
	HallOfFameAt mo.Option[time.Time]

	ViewCount               int64
	ViewMilestoneCount      int64
	ViewMilestoneArchivedAt mo.Option[time.Time]
}

type LoveStats struct {
	Score        model.ScoreType
	FirstLovedAt mo.Option[time.Time]
	LastLovedAt  mo.Option[time.Time]
}
