package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/sqlite/image"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func mapRow[C persist.ImageListCursor](row gen.Image, cursor C) (persist.ImageWithCursor[C], error) {
	image, err := image.MapImage(row)
	if err != nil {
		return persist.ImageWithCursor[C]{}, err
	}

	return persist.ImageWithCursor[C]{
		Image:  image,
		Cursor: cursor,
	}, nil
}

func mapRows[T any, C persist.ImageListCursor](
	rows []T,
	getImage func(T) gen.Image,
	getCursor func(T) C,
) ([]persist.ImageWithCursor[C], error) {
	imagesWithCursor := make([]persist.ImageWithCursor[C], len(rows))

	for i, row := range rows {
		imageWithCursor, err := mapRow(getImage(row), getCursor(row))
		if err != nil {
			return nil, err
		}

		imagesWithCursor[i] = imageWithCursor
	}

	return imagesWithCursor, nil
}
