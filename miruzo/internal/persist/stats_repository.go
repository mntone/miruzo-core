package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type StatsRepository interface {
	ApplyLove(
		ctx context.Context,
		ingestID model.IngestIDType,
		scoreDelta model.ScoreType,
		lovedAt time.Time,
		periodStartAt time.Time,
	) (LoveStats, error)

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
