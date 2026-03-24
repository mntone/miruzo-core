package media

// --- format ---

type ImageFormat uint8

const (
	ImageFormatBitmap ImageFormat = 11
	ImageFormatTIFF   ImageFormat = 12
	ImageFormatGIF    ImageFormat = 13
	ImageFormatPNG    ImageFormat = 14

	ImageFormatJPEG     ImageFormat = 31
	ImageFormatJPEG2000 ImageFormat = 32
	ImageFormatJPEGXR   ImageFormat = 33
	ImageFormatJPEGXL   ImageFormat = 34

	ImageFormatWebP ImageFormat = 61
	ImageFormatAVIF ImageFormat = 62

	ImageFormatAVCI ImageFormat = 91
	ImageFormatHEIF ImageFormat = 92
)

// --- encoding ---

type ImageEncoding uint8

const (
	ImageEncodingUnspecified ImageEncoding = iota
	ImageEncodingJPEG
	ImageEncodingWebP
	ImageEncodingLosslessWebP
)

// --- variant ---

type LayerIDType uint32

const FallbackLayerID LayerIDType = 9

type QualityType uint16

type VariantSpec struct {
	LayerID  LayerIDType
	Width    uint16
	Encoding ImageEncoding
	Quality  QualityType
	Required bool
}

type VariantLayerSpec struct {
	LayerID  LayerIDType
	Variants []VariantSpec
}

type VariantLayersSpec []VariantLayerSpec
