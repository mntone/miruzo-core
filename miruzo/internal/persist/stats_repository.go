package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type StatsRepository interface {
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
