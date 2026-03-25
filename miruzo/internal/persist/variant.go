package persist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/samber/mo"
)

type Variant struct {
	RelativePath string             `json:"rel"`
	LayerID      media.LayerIDType  `json:"layer_id"`
	Format       media.ImageFormat  `json:"format"`
	Codecs       string             `json:"codecs,omitempty"`
	Bytes        uint32             `json:"bytes"`
	Width        uint16             `json:"width"`
	Height       uint16             `json:"height"`
	Quality      *media.QualityType `json:"quality"`
}

func (v *Variant) ToDomain() media.Variant {
	return media.Variant{
		RelativePath: v.RelativePath,
		LayerID:      v.LayerID,
		Format:       v.Format,
		Codecs:       v.Codecs,
		Bytes:        v.Bytes,
		Width:        v.Width,
		Height:       v.Height,
		Quality:      mo.PointerToOption(v.Quality),
	}
}

type Variants []Variant

func (layers Variants) ToDomain() media.Variants {
	variants := make(media.Variants, len(layers))
	for i, variant := range layers {
		variants[i] = variant.ToDomain()
	}

	return variants
}
