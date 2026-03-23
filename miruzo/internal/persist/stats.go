package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Stats struct {
	// IngestID is numeric primary key assigned in the database.
	IngestID model.IngestIDType
	// Score is a user-adjustable ranking score.
	Score            model.ScoreType
	ScoreEvaluated   model.ScoreType
	ScoreEvaluatedAt mo.Option[time.Time]

	// LastViewedAt is the time of the most recent view, or unset if never viewed.
	LastViewedAt mo.Option[time.Time]
	// FirstLovedAt is the time of the first love action, or unset if never loved.
	FirstLovedAt mo.Option[time.Time]
	// LastLovedAt is the time of the latest love action, or unset if never loved.
	LastLovedAt mo.Option[time.Time]
	// HallOfFameAt is the time when the image entered the hall of fame, or unset if it has not.
	HallOfFameAt mo.Option[time.Time]

	// ViewCount is how many times this image has been viewed.
	ViewCount int64
	// ViewMilestoneCount is the highest view milestone reached so far.
	ViewMilestoneCount int64
	// ViewMilestoneArchivedAt is timestamp when the latest view milestone was reached, or unset.
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

// Hall of fame statistics for a single image.
type HallOfFameStats struct {
	// HallOfFameAt is the time when the image entered the hall of fame, or unset if it has not.
	HallOfFameAt mo.Option[time.Time]
}
