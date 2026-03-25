package media

import "github.com/samber/mo"

type Variant struct {
	// RelativePath is the server-relative asset path.
	RelativePath string
	// LayerID identifies which layer this variant belongs to.
	LayerID LayerIDType
	// Format is the image container format.
	Format ImageFormat
	// Codecs is an optional codec hint (for example, vp8 or vp8l).
	Codecs string
	// Bytes is the encoded file size in bytes.
	Bytes uint32
	// Width is the variant width in pixels.
	Width uint16
	// Height is the variant height in pixels.
	Height uint16
	// Quality is the optional encoding quality for lossy outputs.
	Quality mo.Option[QualityType]
}

func (v Variant) IsFallback() bool {
	return v.LayerID == FallbackLayerID
}

type Variants []Variant

type VariantLayers [][]Variant
