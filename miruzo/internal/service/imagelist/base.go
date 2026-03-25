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

type imageListFunc[S model.ImageListCursorScalar, T any] func(
	requestContext context.Context,
	spec T,
) ([]persist.ImageWithCursorKey[S], error)

func listBase[S model.ImageListCursorScalar, T any](
	requestContext context.Context,
	limit uint16,
	excludeFormats []media.ImageFormat,
	loadFn imageListFunc[S, T],
	spec T,
	variantLayersBuilder *media.VariantLayersBuilder,
	retryPolicy backoff.Policy,
) (Result[S], error) {
	imagesWithCursor, err := retry.Retry(
		requestContext,
		retryPolicy,
		func(retryContext context.Context) ([]persist.ImageWithCursorKey[S], error) {
			return loadFn(retryContext, spec)
		},
	)
	if err != nil {
		return Result[S]{}, serviceerror.MapPersistError(err)
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

	var nextCursor mo.Option[model.ImageListCursorKey[S]]
	if hasNext {
		image := imagesWithCursor[n-1]
		nextCursor = mo.Some(model.ImageListCursorKey[S]{
			Primary:   image.PrimaryKey,
			Secondary: image.Image.IngestID,
		})
	}

	result := Result[S]{
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
		CursorKey: params.Cursor,
		MaxCount:  params.Limit + 1,
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
