package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapLatestRows(rows []gen.Image) ([]persist.ImageWithCursor[time.Time], error) {
	imagesWithCursor := make([]persist.ImageWithCursor[time.Time], len(rows))

	for i, row := range rows {
		imageWithCursor, err := mapRow(row, row.IngestedAt.Time)
		if err != nil {
			return nil, err
		}

		imagesWithCursor[i] = imageWithCursor
	}

	return imagesWithCursor, nil
}

func (repo repository) ListLatest(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesLatest(
			ctx,
			int32(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListLatest", err)
		}

		return mapLatestRows(rows)
	}

	rows, err := repo.queries.ListImagesLatestAfter(
		ctx,
		gen.ListImagesLatestAfterParams{
			IngestedAt: shared.PgtypeTimestampFromTime(cursor),
			Limit:      int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListLatest", err)
	}

	return mapLatestRows(rows)
}
