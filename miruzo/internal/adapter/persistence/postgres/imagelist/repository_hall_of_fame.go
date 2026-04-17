package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/dberrors"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapHallOfFameRows(rows []gen.ListImagesHallOfFameRow) ([]persist.ImageWithCursorKey[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesHallOfFameRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesHallOfFameRow) time.Time {
			return *row.HallOfFameAt
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
			return *row.HallOfFameAt
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
			int32(spec.MaxCount),
		)
		if err != nil {
			return nil, dberrors.ToPersist("ListHallOfFame", err)
		}

		return mapHallOfFameRows(rows)
	}

	rows, err := repo.queries.ListImagesHallOfFameAfter(
		ctx,
		gen.ListImagesHallOfFameAfterParams{
			CursorAt: &cursor.Primary,
			CursorID: cursor.Secondary,
			MaxCount: int32(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, dberrors.ToPersist("ListHallOfFame", err)
	}

	return mapHallOfFameAfterRows(rows)
}
