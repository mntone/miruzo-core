package config

import "github.com/mntone/miruzo-core/miruzo/internal/model/media"

type ImageFormat uint8

const (
	ImageFormatBitmap ImageFormat = 11
	ImageFormatTiff   ImageFormat = 12
	ImageFormatGif    ImageFormat = 13
	ImageFormatPng    ImageFormat = 14

	ImageFormatJpeg     ImageFormat = 31
	ImageFormatJpeg2000 ImageFormat = 32
	ImageFormatXr       ImageFormat = 33
	ImageFormatXl       ImageFormat = 34

	ImageFormatWebp ImageFormat = 61
	ImageFormatAvif ImageFormat = 62

	ImageFormatAvci ImageFormat = 91
	ImageFormatHeif ImageFormat = 92
)

type VariantFormat struct {
	Format        ImageFormat
	Codecs        string
	FileExtension string
}

type VariantConfig struct {
	LayerID  media.LayerIDType
	Width    uint16
	Format   VariantFormat
	Quality  media.QualityType
	Required bool
}

type VariantLayerConfig struct {
	LayerID  media.LayerIDType
	Variants []VariantConfig
}

var jpegVariantFormat = VariantFormat{
	Format:        ImageFormatJpeg,
	Codecs:        "",
	FileExtension: ".jpg",
}

var webpVariantFormat = VariantFormat{
	Format:        ImageFormatWebp,
	Codecs:        "vp8",
	FileExtension: ".webp",
}

var losslessWebpVariantFormat = VariantFormat{
	Format:        ImageFormatWebp,
	Codecs:        "vp8l",
	FileExtension: ".webp",
}

func DefaultVariantLayerConfig() []VariantLayerConfig {
	return []VariantLayerConfig{
		{
			LayerID: 1,
			Variants: []VariantConfig{
				{
					LayerID:  1,
					Width:    320,
					Format:   webpVariantFormat,
					Quality:  80,
					Required: true,
				},
				{
					LayerID:  1,
					Width:    480,
					Format:   webpVariantFormat,
					Quality:  70,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    640,
					Format:   webpVariantFormat,
					Quality:  60,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    960,
					Format:   webpVariantFormat,
					Quality:  50,
					Required: false,
				},
				{
					LayerID:  1,
					Width:    1120,
					Format:   webpVariantFormat,
					Quality:  40,
					Required: false,
				},
			},
		},
		{
			LayerID: media.FallbackLayerID,
			Variants: []VariantConfig{
				{
					LayerID:  media.FallbackLayerID,
					Width:    320,
					Format:   jpegVariantFormat,
					Quality:  85,
					Required: true,
				},
			},
		},
	}
}
