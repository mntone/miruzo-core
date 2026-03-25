package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func (srv Service) ListLatest(
	requestContext context.Context,
	params *Params[time.Time],
) (Result[time.Time], error) {
	return list(
		requestContext,
		srv.repository.ListLatest,
		params,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}

func (srv Service) ListChronological(
	requestContext context.Context,
	params *Params[time.Time],
) (Result[time.Time], error) {
	return list(
		requestContext,
		srv.repository.ListChronological,
		params,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}

func (srv Service) ListRecently(
	requestContext context.Context,
	params *Params[time.Time],
) (Result[time.Time], error) {
	return list(
		requestContext,
		srv.repository.ListRecently,
		params,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}

func (srv Service) ListFirstLove(
	requestContext context.Context,
	params *Params[time.Time],
) (Result[time.Time], error) {
	return list(
		requestContext,
		srv.repository.ListFirstLove,
		params,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}

func (srv Service) ListHallOfFame(
	requestContext context.Context,
	params *Params[time.Time],
) (Result[time.Time], error) {
	return list(
		requestContext,
		srv.repository.ListHallOfFame,
		params,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}

func (srv Service) ListEngaged(
	requestContext context.Context,
	params *Params[model.ScoreType],
) (Result[model.ScoreType], error) {
	spec := persist.EngagedImageListSpec{
		ScoreThreshold: srv.engagedScoreThreshold,
		ImageListSpec: persist.ImageListSpec[model.ScoreType]{
			CursorKey: params.Cursor,
			MaxCount:  params.Limit + 1,
		},
	}
	return listBase(
		requestContext,
		params.Limit,
		params.ExcludeFormats,
		srv.repository.ListEngaged,
		spec,
		srv.variantLayersBuilder,
		srv.backoff,
	)
}
