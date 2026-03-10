package persist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type ImageListSpec[C ImageListCursor] struct {
	Cursor mo.Option[C]
	Limit  uint16
}

type EngagedImageListSpec struct {
	ImageListSpec[model.ScoreType]
	ScoreThreshold model.ScoreType
}

// ImageListRepository provides read operations for image list endpoints.
// All methods use cursor-based pagination and return items ordered by each list kind.
type ImageListRepository interface {
	// ListLatest returns images ordered by ingested_at DESC.
	ListLatest(
		requestContext context.Context,
		spec ImageListSpec[time.Time],
	) ([]ImageWithCursor[time.Time], error)

	// ListChronological returns images ordered by captured_at DESC.
	ListChronological(
		requestContext context.Context,
		spec ImageListSpec[time.Time],
	) ([]ImageWithCursor[time.Time], error)

	// ListRecently returns images ordered by last_viewed_at DESC.
	ListRecently(
		requestContext context.Context,
		spec ImageListSpec[time.Time],
	) ([]ImageWithCursor[time.Time], error)

	// ListFirstLove returns images ordered by first_loved_at DESC.
	ListFirstLove(
		requestContext context.Context,
		spec ImageListSpec[time.Time],
	) ([]ImageWithCursor[time.Time], error)

	// ListHallOfFame returns images ordered by hall_of_fame_at DESC.
	ListHallOfFame(
		requestContext context.Context,
		spec ImageListSpec[time.Time],
	) ([]ImageWithCursor[time.Time], error)

	// ListEngaged returns images ordered by score_evaluated DESC.
	ListEngaged(
		requestContext context.Context,
		spec EngagedImageListSpec,
	) ([]ImageWithCursor[model.ScoreType], error)
}
