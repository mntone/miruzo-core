package domain

import "github.com/mntone/miruzo-core/miruzo/internal/domain/media"

func NewTestVariantLayersBuilder() *media.VariantLayersBuilder {
	return media.NewVariantLayerBuilder(media.VariantLayersSpec{
		{
			LayerID: 1,
			Variants: []media.VariantSpec{
				{
					LayerID:  1,
					Width:    320,
					Encoding: media.ImageEncodingWebP,
					Quality:  80,
					Required: true,
				},
			},
		},
	})
}
