package testutil

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/samber/mo"
)

type TestRepository interface {
	CreateIngest(
		ctx context.Context,
		id int64,
		relativePath string,
		fingerprint string,
		ingestedAt time.Time,
		capturedAt time.Time,
	) error

	CreateImage(
		ctx context.Context,
		id int64,
		ingestedAt time.Time,
		original media.Variant,
		fallback mo.Option[media.Variant],
		variants []media.Variant,
	) error

	CreateStat(
		ctx context.Context,
		id int64,
		score int16,
		scoreEvaluated int16,
		lastViewedAt mo.Option[time.Time],
		firstLovedAt mo.Option[time.Time],
		lastLovedAt mo.Option[time.Time],
		hallOfFameAt mo.Option[time.Time],
		viewCount int64,
	) error
}
