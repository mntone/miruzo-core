package imagelist

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/mysql/dberrors"
	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapFirstLoveRows(rows []gen.ListImagesFirstLoveRow) ([]persist.ImageWithCursorKey[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesFirstLoveRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesFirstLoveRow) time.Time {
			return persistshared.TimeFromSql(row.FirstLovedAt)
		},
	)
}

func mapFirstLoveAfterRows(rows []gen.ListImagesFirstLoveAfterRow) ([]persist.ImageWithCursorKey[time.Time], error) {
	return mapRows(
		rows,
		func(row gen.ListImagesFirstLoveAfterRow) gen.Image {
			return row.Image
		},
		func(row gen.ListImagesFirstLoveAfterRow) time.Time {
			return persistshared.TimeFromSql(row.FirstLovedAt)
		},
	)
}

func (repo repository) ListFirstLove(
	ctx context.Context,
	spec persist.ImageListSpec[time.Time],
) ([]persist.ImageWithCursorKey[time.Time], error) {
	cursor, present := spec.CursorKey.Get()
	if !present {
		rows, err := repo.queries.ListImagesFirstLove(
			ctx,
			int32(spec.MaxCount),
		)
		if err != nil {
			return nil, dberrors.ToPersist("ListFirstLove", err)
		}

		return mapFirstLoveRows(rows)
	}

	rows, err := repo.queries.ListImagesFirstLoveAfter(
		ctx,
		gen.ListImagesFirstLoveAfterParams{
			CursorAt: persistshared.NullTimeFromTime(cursor.Primary),
			CursorID: cursor.Secondary,
			Limit:    int32(spec.MaxCount),
		},
	)
	if err != nil {
		return nil, dberrors.ToPersist("ListFirstLove", err)
	}

	return mapFirstLoveAfterRows(rows)
}
