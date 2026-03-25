package media

import (
	"encoding/json"
	"errors"
)

// --- format ---

type ImageFormat uint8

const (
	ImageFormatUnspecified ImageFormat = 0

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

const (
	imageFormatStringBitmap = "bmp"
	imageFormatStringTIFF   = "tiff"
	imageFormatStringGIF    = "gif"
	imageFormatStringPNG    = "png"

	imageFormatStringJPEG     = "jpeg"
	imageFormatStringJPEG2000 = "jp2"
	imageFormatStringJPEGXR   = "jxr"
	imageFormatStringJPEGXL   = "jxl"

	imageFormatStringWebP = "webp"
	imageFormatStringAVIF = "avif"

	imageFormatStringAVCI = "avci"
	imageFormatStringHEIF = "heif"
)

func ParseImageFormat(imageFormatString string) (ImageFormat, bool) {
	switch imageFormatString {
	case imageFormatStringBitmap:
		return ImageFormatBitmap, true
	case imageFormatStringTIFF:
		return ImageFormatTIFF, true
	case imageFormatStringGIF:
		return ImageFormatGIF, true
	case imageFormatStringPNG:
		return ImageFormatPNG, true

	case imageFormatStringJPEG:
		return ImageFormatJPEG, true
	case imageFormatStringJPEG2000:
		return ImageFormatJPEG2000, true
	case imageFormatStringJPEGXR:
		return ImageFormatJPEGXR, true
	case imageFormatStringJPEGXL:
		return ImageFormatJPEGXL, true

	case imageFormatStringWebP:
		return ImageFormatWebP, true
	case imageFormatStringAVIF:
		return ImageFormatAVIF, true

	case imageFormatStringAVCI:
		return ImageFormatAVCI, true
	case imageFormatStringHEIF:
		return ImageFormatHEIF, true
	}

	return ImageFormatUnspecified, false
}

func (v ImageFormat) String() string {
	switch v {
	case ImageFormatBitmap:
		return imageFormatStringBitmap
	case ImageFormatTIFF:
		return imageFormatStringTIFF
	case ImageFormatGIF:
		return imageFormatStringGIF
	case ImageFormatPNG:
		return imageFormatStringPNG

	case ImageFormatJPEG:
		return imageFormatStringJPEG
	case ImageFormatJPEG2000:
		return imageFormatStringJPEG2000
	case ImageFormatJPEGXR:
		return imageFormatStringJPEGXR
	case ImageFormatJPEGXL:
		return imageFormatStringJPEGXL

	case ImageFormatWebP:
		return imageFormatStringWebP
	case ImageFormatAVIF:
		return imageFormatStringAVIF

	case ImageFormatAVCI:
		return imageFormatStringAVCI
	case ImageFormatHEIF:
		return imageFormatStringHEIF
	}

	return "bin"
}

func (v ImageFormat) MarshalJSON() ([]byte, error) {
	if v == ImageFormatUnspecified {
		return nil, errors.New("invalid image format")
	}

	return json.Marshal(v.String())
}

func (v *ImageFormat) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsedValue, ok := ParseImageFormat(s)
	if !ok {
		return errors.New("unknown image format")
	}

	*v = parsedValue
	return nil
}

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
