package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRecentlyRows(rows []gen.ListImagesRecentlyRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesRecentlyRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesRecentlyRow) time.Time {
			return shared.TimeFromSql(row.LastViewedAt)
		},
	)
}

func mapRecentlyAfterRows(rows []gen.ListImagesRecentlyAfterRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesRecentlyAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesRecentlyAfterRow) time.Time {
			return shared.TimeFromSql(row.LastViewedAt)
		},
	)
}

func (repo repository) ListRecently(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesRecently(
			ctx,
			int64(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapSQLiteError("ListRecently", err)
		}

		return mapRecentlyRows(rows)
	}

	rows, err := repo.queries.ListImagesRecentlyAfter(
		ctx,
		gen.ListImagesRecentlyAfterParams{
			LastViewedAt: shared.NullTimeFromTime(cursor),
			Limit:        int64(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapSQLiteError("ListRecently", err)
	}

	return mapRecentlyAfterRows(rows)
}
