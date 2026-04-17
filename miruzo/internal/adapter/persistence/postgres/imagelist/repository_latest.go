package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapLatestRows(rows []gen.Image) ([]persist.ImageWithCursorKey[time.Time], error) {
	imagesWithCursor := make([]persist.ImageWithCursorKey[time.Time], len(rows))

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
) ([]persist.ImageWithCursorKey[time.Time], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesLatest(
			ctx,
			int32(spec.MaxCount),
		)
		if err != nil {
			return nil, dberrors.ToPersist("ListLatest", err)
		}

		return mapLatestRows(rows)
	}

	rows, err := repo.queries.ListImagesLatestAfter(
		ctx,
		gen.ListImagesLatestAfterParams{
			CursorAt: cursor.Primary,
			CursorID: cursor.Secondary,
			MaxCount: int32(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, dberrors.ToPersist("ListLatest", err)
	}

	return mapLatestRows(rows)
}
