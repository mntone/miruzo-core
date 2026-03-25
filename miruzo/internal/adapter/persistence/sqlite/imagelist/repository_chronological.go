package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapChronologicalRows(rows []gen.ListImagesChronologicalRow) ([]persist.ImageWithCursorKey[time.Time], error) {
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

func mapChronologicalAfterRows(rows []gen.ListImagesChronologicalAfterRow) ([]persist.ImageWithCursorKey[time.Time], error) {
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
) ([]persist.ImageWithCursorKey[time.Time], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesChronological(
			ctx,
			int64(spec.MaxCount),
		)
		if err != nil {
			return nil, shared.MapSQLiteError("ListChronological", err)
		}

		return mapChronologicalRows(rows)
	}

	rows, err := repo.queries.ListImagesChronologicalAfter(
		ctx,
		gen.ListImagesChronologicalAfterParams{
			CursorAt: cursor.Primary,
			CursorID: cursor.Secondary,
			MaxCount: int64(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, shared.MapSQLiteError("ListChronological", err)
	}

	return mapChronologicalAfterRows(rows)
}
