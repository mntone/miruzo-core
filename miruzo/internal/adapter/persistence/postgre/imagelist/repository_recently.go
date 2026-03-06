package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRecentlyRows(rows []gen.ListImagesRecentlyRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesRecentlyRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesRecentlyRow) time.Time {
			return shared.TimeFromPgtype(row.LastViewedAt)
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
			return shared.TimeFromPgtype(row.LastViewedAt)
		},
	)
}

func (repo *repository) ListRecently(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesRecently(
			ctx,
			int32(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListRecently", err)
		}

		return mapRecentlyRows(rows)
	}

	rows, err := repo.queries.ListImagesRecentlyAfter(
		ctx,
		gen.ListImagesRecentlyAfterParams{
			LastViewedAt: shared.PgtypeTimestampFromTime(cursor),
			Limit:        int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListRecently", err)
	}

	return mapRecentlyAfterRows(rows)
}
