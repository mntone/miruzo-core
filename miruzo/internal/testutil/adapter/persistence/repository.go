package persistence

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type TestRepository interface {
	CreateIngest(
		ctx context.Context,
		id model.IngestIDType,
		relativePath string,
		fingerprint string,
		ingestedAt time.Time,
		capturedAt time.Time,
	) error

	CreateImage(
		ctx context.Context,
		id model.IngestIDType,
		ingestedAt time.Time,
		original persist.Variant,
		fallback mo.Option[persist.Variant],
		variants []persist.Variant,
	) error

	CreateStat(
		ctx context.Context,
		id model.IngestIDType,
		score model.ScoreType,
		scoreEvaluated model.ScoreType,
		lastViewedAt mo.Option[time.Time],
		firstLovedAt mo.Option[time.Time],
		lastLovedAt mo.Option[time.Time],
		hallOfFameAt mo.Option[time.Time],
		viewCount int64,
	) error

	ExecuteStatement(ctx context.Context, stmt string, delete bool) error
	ExecuteStatementAndReturnRowCount(ctx context.Context, stmt string, delete bool) (int64, error)

	TruncateActions(ctx context.Context) error
	TruncateJobs(ctx context.Context) error
	TruncateStats(ctx context.Context) error
}
