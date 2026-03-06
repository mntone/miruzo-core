package persist

import (
	"time"

	"github.com/samber/mo"
)

type Stat struct {
	IngestID         IngestID
	Score            int16
	ScoreEvaluated   int16
	ScoreEvaluatedAt mo.Option[time.Time]

	LastViewedAt mo.Option[time.Time]
	FirstLovedAt mo.Option[time.Time]
	LastLovedAt  mo.Option[time.Time]
	HallOfFameAt mo.Option[time.Time]

	ViewCount               int64
	ViewMilestoneCount      int64
	ViewMilestoneArchivedAt mo.Option[time.Time]
}
