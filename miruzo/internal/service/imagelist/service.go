package imagelist

import (
	"context"
	"time"

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
	params *Params[int16],
) (Result[int16], error) {
	spec := persist.EngagedImageListSpec{
		ImageListSpec: persist.ImageListSpec[int16]{
			Cursor: params.Cursor,
			Limit:  params.Limit + 1,
		},
		ScoreThreshold: srv.engagedScoreThreshold,
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
