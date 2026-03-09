package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapLatestRows(rows []gen.Image) ([]persist.ImageWithCursor[time.Time], error) {
	imagesWithCursor := make([]persist.ImageWithCursor[time.Time], len(rows))

	for i, row := range rows {
		imageWithCursor, err := mapRow(row, row.IngestedAt)
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
			int64(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapSQLiteError("ListLatest", err)
		}

		return mapLatestRows(rows)
	}

	rows, err := repo.queries.ListImagesLatestAfter(
		ctx,
		gen.ListImagesLatestAfterParams{
			IngestedAt: cursor,
			Limit:      int64(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapSQLiteError("ListLatest", err)
	}

	return mapLatestRows(rows)
}
