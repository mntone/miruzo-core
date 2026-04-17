package imagelist

import (
	"context"
	"time"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapHallOfFameRows(rows []gen.ListImagesHallOfFameRow) ([]persist.ImageWithCursorKey[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesHallOfFameRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesHallOfFameRow) time.Time {
			return persistshared.TimeFromSql(row.HallOfFameAt)
		},
	)
}

func mapHallOfFameAfterRows(rows []gen.ListImagesHallOfFameAfterRow) ([]persist.ImageWithCursorKey[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesHallOfFameAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesHallOfFameAfterRow) time.Time {
			return persistshared.TimeFromSql(row.HallOfFameAt)
		},
	)
}

func (repo repository) ListHallOfFame(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursorKey[time.Time], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesHallOfFame(
			ctx,
			int64(spec.MaxCount),
		)
		if err != nil {
			return nil, dberrors.ToPersist("ListHallOfFame", err)
		}

		return mapHallOfFameRows(rows)
	}

	rows, err := repo.queries.ListImagesHallOfFameAfter(
		ctx,
		gen.ListImagesHallOfFameAfterParams{
			CursorAt: persistshared.NullTimeFromTime(cursor.Primary),
			CursorID: cursor.Secondary,
			MaxCount: int64(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, dberrors.ToPersist("ListHallOfFame", err)
	}

	return mapHallOfFameAfterRows(rows)
}
