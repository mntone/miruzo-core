package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgre/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapHallOfFameRows(rows []gen.ListImagesHallOfFameRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesHallOfFameRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesHallOfFameRow) time.Time {
			return shared.TimeFromPgtype(row.HallOfFameAt)
		},
	)
}

func mapHallOfFameAfterRows(rows []gen.ListImagesHallOfFameAfterRow) ([]persist.ImageWithCursor[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesHallOfFameAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesHallOfFameAfterRow) time.Time {
			return shared.TimeFromPgtype(row.HallOfFameAt)
		},
	)
}

func (repo *repository) ListHallOfFame(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursor[time.Time], error) {
	cursor, present := spec.Cursor.Get()
	if !present {
		rows, err := repo.queries.ListImagesHallOfFame(
			ctx,
			int32(spec.Limit),
		)
		if err != nil {
			return nil, shared.MapPostgreError("ListHallOfFame", err)
		}

		return mapHallOfFameRows(rows)
	}

	rows, err := repo.queries.ListImagesHallOfFameAfter(
		ctx,
		gen.ListImagesHallOfFameAfterParams{
			HallOfFameAt: shared.PgtypeTimestampFromTime(cursor),
			Limit:        int32(spec.Limit),
		},
	)
	if err != nil {
		return nil, shared.MapPostgreError("ListHallOfFame", err)
	}

	return mapHallOfFameAfterRows(rows)
}
