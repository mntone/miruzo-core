package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRecentlyRows(rows []gen.ListImagesRecentlyRow) ([]persist.ImageWithCursorKey[time.Time], error) {
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

func mapRecentlyAfterRows(rows []gen.ListImagesRecentlyAfterRow) ([]persist.ImageWithCursorKey[time.Time], error) {
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
) ([]persist.ImageWithCursorKey[time.Time], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesRecently(
			ctx,
			int32(spec.MaxCount),
		)
		if err != nil {
			return nil, dberrors.ToPersist("ListRecently", err)
		}

		return mapRecentlyRows(rows)
	}

	rows, err := repo.queries.ListImagesRecentlyAfter(
		ctx,
		gen.ListImagesRecentlyAfterParams{
			CursorAt: &cursor.Primary,
			CursorID: cursor.Secondary,
			MaxCount: int32(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, dberrors.ToPersist("ListRecently", err)
	}

	return mapRecentlyAfterRows(rows)
}
