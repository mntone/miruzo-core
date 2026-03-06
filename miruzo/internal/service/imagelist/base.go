package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/samber/mo"
)

func mapPageBounds(total int, limit int) (n int, hasNext bool) {
	n = total

	if n > limit {
		n = limit
		hasNext = true
	}

	return
}

type imageListFunc[C persist.ImageListCursor, S any] func(
	requestContext context.Context,
	spec S,
) ([]persist.ImageWithCursor[C], error)

func listBase[C persist.ImageListCursor, S any](
	requestContext context.Context,
	loadFn imageListFunc[C, S],
	params *Params[C],
	spec S,
	retryPolicy backoff.Policy,
) (Result[C], error) {
	imagesWithCursor, err := retry.Retry(
		requestContext,
		retryPolicy,
		func(retryContext context.Context) ([]persist.ImageWithCursor[C], error) {
			return loadFn(retryContext, spec)
		},
	)
	if err != nil {
		return Result[C]{}, serviceerror.MapPersistError(err)
	}

	n, hasNext := mapPageBounds(len(imagesWithCursor), int(params.Limit))

	images := make([]persist.Image, n)
	for i := range n {
		images[i] = imagesWithCursor[i].Image
	}

	var nextCursor mo.Option[C]
	if hasNext {
		nextCursor = mo.Some(imagesWithCursor[n-1].Cursor)
	}

	result := Result[C]{
		Items:  images,
		Cursor: nextCursor,
	}
	return result, nil
}

func list(
	requestContext context.Context,
	loadFn imageListFunc[time.Time, persist.ImageListSpec[time.Time]],
	params *Params[time.Time],
	retryPolicy backoff.Policy,
) (Result[time.Time], error) {
	spec := persist.ImageListSpec[time.Time]{
		Cursor: params.Cursor,
		Limit:  params.Limit + 1,
	}
	return listBase(
		requestContext,
		loadFn,
		params,
		spec,
		retryPolicy,
	)
}
