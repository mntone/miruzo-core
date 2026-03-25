package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
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
	limit uint16,
	excludeFormats []media.ImageFormat,
	loadFn imageListFunc[C, S],
	spec S,
	variantLayersBuilder *media.VariantLayersBuilder,
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

	n, hasNext := mapPageBounds(len(imagesWithCursor), int(limit))

	options := media.VariantFilterOptions{
		IncludeFormatSet: media.ComputeAllowedFormats(excludeFormats),
		KeepFallback:     true,
	}
	images := make([]model.Image, n)
	for i := range n {
		variants := imagesWithCursor[i].Image.Layers.ToDomain().FilterWith(options)
		layers := variantLayersBuilder.GroupVariantsByLayer(variants)
		images[i] = imagesWithCursor[i].Image.ToDTO(layers)
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
	variantLayersBuilder *media.VariantLayersBuilder,
	retryPolicy backoff.Policy,
) (Result[time.Time], error) {
	spec := persist.ImageListSpec[time.Time]{
		Cursor: params.Cursor,
		Limit:  params.Limit + 1,
	}
	return listBase(
		requestContext,
		params.Limit,
		params.ExcludeFormats,
		loadFn,
		spec,
		variantLayersBuilder,
		retryPolicy,
	)
}
