package imagelist

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/sqlite/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
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

	original, err := shared.MapVariant(row.Original)
	if err != nil {
		return persist.ImageWithCursor[C]{}, err
	}

	fallback, err := shared.MapNullableVariant(row.Fallback)
	if err != nil {
		return persist.ImageWithCursor[C]{}, err
	}

	variants, err := shared.MapVariants(row.Variants)
	if err != nil {
		return persist.ImageWithCursor[C]{}, err
	}

	return persist.ImageWithCursor[C]{
		Image: persist.Image{
			IngestID:   row.IngestID,
			IngestedAt: row.IngestedAt,
			Type:       imageType,
			Original:   original,
			Fallback:   fallback,
			Variants:   variants,
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
