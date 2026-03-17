package image

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/database/postgres/gen"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
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

	return persist.Image{
		IngestID:   r.IngestID,
		IngestedAt: r.IngestedAt,
		Type:       imageType,
		Original:   r.Original,
		Fallback:   mo.PointerToOption(r.Fallback),
		Variants:   r.Variants,
	}, nil
}
