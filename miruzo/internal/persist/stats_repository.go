package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type StatsRepository interface {
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
		periodStartAt time.Time,
		dayStartOffset time.Duration,
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
}
