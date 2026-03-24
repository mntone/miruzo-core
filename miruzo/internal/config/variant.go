package config

import "github.com/mntone/miruzo-core/miruzo/internal/domain/media"

func DefaultVariantLayerConfig() media.VariantLayersSpec {
	return media.VariantLayersSpec{
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
				{
					LayerID:  1,
					Width:    480,
					Encoding: media.ImageEncodingWebP,
					Quality:  70,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    640,
					Encoding: media.ImageEncodingWebP,
					Quality:  60,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    960,
					Encoding: media.ImageEncodingWebP,
					Quality:  50,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    1120,
					Encoding: media.ImageEncodingWebP,
					Quality:  40,
					Required: false,
				},
			},
		},
		{
			LayerID: media.FallbackLayerID,
			Variants: []media.VariantSpec{
				{
					LayerID:  media.FallbackLayerID,
					Width:    320,
					Encoding: media.ImageEncodingJPEG,
					Quality:  85,
					Required: true,
				},
			},
		},
	}
}
