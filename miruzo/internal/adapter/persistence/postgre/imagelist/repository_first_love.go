package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapFirstLoveRows(rows []gen.ListImagesFirstLoveRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesFirstLoveRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesFirstLoveRow) time.Time {
			return shared.TimeFromPgtype(row.FirstLovedAt)
		},
	)
}

func mapFirstLoveAfterRows(rows []gen.ListImagesFirstLoveAfterRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesFirstLoveAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesFirstLoveAfterRow) time.Time {
			return shared.TimeFromPgtype(row.FirstLovedAt)
		},
	)
}

func (repo *repository) ListFirstLove(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesFirstLove(
			ctx,
			int32(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListFirstLove", err)
		}

		return mapFirstLoveRows(rows)
	}

	rows, err := repo.queries.ListImagesFirstLoveAfter(
		ctx,
		gen.ListImagesFirstLoveAfterParams{
			FirstLovedAt: shared.PgtypeTimestampFromTime(cursor),
			Limit:        int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListFirstLove", err)
	}

	return mapFirstLoveAfterRows(rows)
}
