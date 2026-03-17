package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapChronologicalRows(rows []gen.ListImagesChronologicalRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesChronologicalRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesChronologicalRow) time.Time {
			return row.CapturedAt
		},
	)
}

func mapChronologicalAfterRows(rows []gen.ListImagesChronologicalAfterRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesChronologicalAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesChronologicalAfterRow) time.Time {
			return row.CapturedAt
		},
	)
}

func (repo repository) ListChronological(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesChronological(
			ctx,
			int32(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListChronological", err)
		}

		return mapChronologicalRows(rows)
	}

	rows, err := repo.queries.ListImagesChronologicalAfter(
		ctx,
		gen.ListImagesChronologicalAfterParams{
			CapturedAt: cursor,
			Limit:      int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListChronological", err)
	}

	return mapChronologicalAfterRows(rows)
}
