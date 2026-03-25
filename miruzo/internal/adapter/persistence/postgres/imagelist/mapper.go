package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/postgres/image"
	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRow[S model.ImageListCursorScalar](row gen.Image, cursor S) (persist.ImageWithCursorKey[S], error) {
	image, err := image.MapImage(row)
	if err != nil {
		return persist.ImageWithCursorKey[S]{}, err
	}

	return persist.ImageWithCursorKey[S]{
		Image:      image,
		PrimaryKey: cursor,
	}, nil
}

func mapRows[T any, S model.ImageListCursorScalar](
	rows []T,
	getImage func(T) gen.Image,
	getCursor func(T) S,
) ([]persist.ImageWithCursorKey[S], error) {
	imagesWithCursor := make([]persist.ImageWithCursorKey[S], len(rows))

	for i, row := range rows {
		imageWithCursor, err := mapRow(getImage(row), getCursor(row))
		if err != nil {
			return nil, err
		}

		imagesWithCursor[i] = imageWithCursor
	}

	return imagesWithCursor, nil
}
