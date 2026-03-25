package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Image struct {
	IngestID   model.IngestIDType
	IngestedAt time.Time
	Type       model.ImageType

	Original Variant
	Fallback mo.Option[Variant]
	Layers   Variants
}

func (e *Image) ToDTO(layers media.VariantLayers) model.Image {
	var fallback mo.Option[media.Variant]
	if f, present := e.Fallback.Get(); present {
		fallback = mo.Some(f.ToDomain())
	}

	return model.Image{
		IngestID:   e.IngestID,
		IngestedAt: e.IngestedAt,
		Type:       e.Type,

		VariantBundle: model.VariantBundle{
			Original: e.Original.ToDomain(),
			Fallback: fallback,
			Layers:   layers,
		},
	}
}

type ImageWithCursorKey[ScalarType model.ImageListCursorScalar] struct {
	Image      Image
	PrimaryKey ScalarType
}

type ImageWithStats struct {
	Image
	Stats model.Stats
}

func (e *ImageWithStats) ToDTO(layers media.VariantLayers) model.ImageWithStats {
	return model.ImageWithStats{
		Image: e.Image.ToDTO(layers),
		Stats: e.Stats,
	}
}
