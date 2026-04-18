package image

import (
	"fmt"

	persistshared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/database/mysql/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func MapImage(r gen.Image) (persist.Image, error) {
	imageType, err := model.ValidateImageType(r.Kind)
	if err != nil {
		return persist.Image{}, fmt.Errorf(
			"%w: ingest_id=%d kind=%d",
			err,
			r.IngestID,
			r.Kind,
		)
	}

	original, err := persistshared.MapVariant(r.Original)
	if err != nil {
		return persist.Image{}, err
	}

	fallback, err := persistshared.MapNullableVariant(r.Fallback)
	if err != nil {
		return persist.Image{}, err
	}

	variants, err := persistshared.MapVariants(r.Variants)
	if err != nil {
		return persist.Image{}, err
	}

	return persist.Image{
		IngestID:   r.IngestID,
		IngestedAt: r.IngestedAt,
		Type:       imageType,
		Original:   original,
		Fallback:   fallback,
		Layers:     variants,
	}, nil
}
