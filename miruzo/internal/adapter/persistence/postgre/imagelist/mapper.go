package imagelist

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/database/postgre/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

func mapRow[C persist.ImageListCursor](row gen.Image, cursor C) (persist.ImageWithCursor[C], error) {
	imageType, err := persist.ParseImageType(row.Kind)
	if err != nil {
		return persist.ImageWithCursor[C]{}, fmt.Errorf(
			"%w: ingest_id=%d kind=%d",
			err,
			row.IngestID,
			row.Kind,
		)
	}

	return persist.ImageWithCursor[C]{
		Image: persist.Image{
			IngestID:   row.IngestID,
			IngestedAt: row.IngestedAt.Time,
			Type:       imageType,
			Original:   row.Original,
			Fallback:   mo.PointerToOption(row.Fallback),
			Variants:   row.Variants,
		},
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
