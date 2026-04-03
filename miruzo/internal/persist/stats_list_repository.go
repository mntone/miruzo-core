package persist

import (
	"context"
	"iter"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type DailyDecayStats struct {
	// IngestID is numeric primary key assigned in the database.
	IngestID model.IngestIDType
	// Score is a user-adjustable ranking score.
	Score model.ScoreType
	// LastViewedAt is the time of the most recent view, or unset if never viewed.
	LastViewedAt mo.Option[time.Time]
}

type StatsListRepository interface {
	IterateStatsForDailyDecay(
		ctx context.Context,
		batchCount int32,
	) iter.Seq2[DailyDecayStats, error]
}
