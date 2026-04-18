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

type StatsRepository interface {
	ApplyDecay(
		ctx context.Context,
		ingestID model.IngestIDType,
		score model.ScoreType,
		evaluatedAt time.Time,
	) error

	ApplyHallOfFameGranted(
		ctx context.Context,
		ingestID model.IngestIDType,
		hallOfFameAt time.Time,
		hallOfFameScoreThreshold model.ScoreType,
	) error

	ApplyHallOfFameRevoked(
		ctx context.Context,
		ingestID model.IngestIDType,
	) error

	ApplyLove(
		ctx context.Context,
		ingestID model.IngestIDType,
		scoreDelta model.ScoreType,
		lovedAt time.Time,
		loveScoreThreshold model.ScoreType,
		periodStartAt time.Time,
	) (model.LoveStats, error)

	ApplyLoveCanceled(
		ctx context.Context,
		ingestID model.IngestIDType,
		scoreDelta model.ScoreType,
		loveCanceledAt time.Time,
		periodStartAt time.Time,
	) (model.LoveStats, error)

	ApplyView(
		ctx context.Context,
		ingestID model.IngestIDType,
		scoreDelta model.ScoreType,
		viewedAt time.Time,
	) error

	ApplyViewWithMilestone(
		ctx context.Context,
		ingestID model.IngestIDType,
		scoreDelta model.ScoreType,
		viewedAt time.Time,
	) error

	IterateStatsForDailyDecay(
		ctx context.Context,
		batchCount int32,
	) iter.Seq2[DailyDecayStats, error]
}
