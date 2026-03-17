package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRecentlyRows(rows []gen.ListImagesRecentlyRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesRecentlyRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesRecentlyRow) time.Time {
			return *row.LastViewedAt
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
			return *row.LastViewedAt
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
			LastViewedAt: &cursor,
			Limit:        int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListRecently", err)
	}

	return mapRecentlyAfterRows(rows)
}
