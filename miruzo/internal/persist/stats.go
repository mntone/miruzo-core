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

// Love statistics for a single image.
type LoveStats struct {
	// Score is a user-adjustable ranking score.
	Score model.ScoreType
	// FirstLovedAt is the time of the first love action, or unset if never loved.
	FirstLovedAt mo.Option[time.Time]
	// LastLovedAt is the time of the latest love action, or unset if never loved.
	LastLovedAt mo.Option[time.Time]
}
