package media

import (
	"fmt"
	"math"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

func TestComputeAllowedFormatsReturnsFallbackFormatsForNilInput(t *testing.T) {
	assert.EqualMap(t, "ComputeAllowedFormats(nil)", ComputeAllowedFormats(nil), defaultAllowedFormats)
}

func TestComputeAllowedFormatsExcludesRequestedFormats(t *testing.T) {
	allowed := ComputeAllowedFormats([]ImageFormat{ImageFormatGIF, ImageFormatWebP})
	assert.EqualMap(t, "ComputeAllowedFormats(allowed)", allowed, map[ImageFormat]struct{}{
		ImageFormatJPEG: {},
		ImageFormatPNG:  {},
	})
}

func TestComputeAllowedFormatsIgnoresUnsupportedFormats(t *testing.T) {
	allowed := ComputeAllowedFormats([]ImageFormat{ImageFormatWebP, ImageFormatJPEGXL})
	assert.EqualMap(t, "ComputeAllowedFormats(allowed)", allowed, map[ImageFormat]struct{}{
		ImageFormatGIF:  {},
		ImageFormatJPEG: {},
		ImageFormatPNG:  {},
	})
}

func makeVariant(
	layerID LayerIDType,
	format ImageFormat,
	width uint16,
) Variant {
	var codecs string
	if format == ImageFormatWebP {
		codecs = "vp8"
	}

	return Variant{
		RelativePath: fmt.Sprintf("l%dw%d/1.%s", layerID, width, format.String()),
		LayerID:      layerID,
		Format:       format,
		Codecs:       codecs,
		Bytes:        234,
		Width:        width,
		Height:       uint16(math.Round(0.75 * float64(width))),
		Quality:      mo.None[QualityType](),
	}
}

func TestFilterWithKeepFallbackTruePreservesFallbackVariant(t *testing.T) {
	variants := Variants{
		makeVariant(1, ImageFormatWebP, 320),
		makeVariant(1, ImageFormatWebP, 480),
		makeVariant(1, ImageFormatJPEGXL, 640),
		makeVariant(1, ImageFormatJPEGXL, 960),
		makeVariant(2, ImageFormatWebP, 640),
		makeVariant(2, ImageFormatWebP, 960),
		makeVariant(FallbackLayerID, ImageFormatJPEG, 320),
	}
	filteredVariants := variants.FilterWith(VariantFilterOptions{
		IncludeFormatSet: ImageFormatSet{
			ImageFormatWebP: {},
		},
		KeepFallback: true,
	})

	assert.LenIs(t, "FilterWith()", filteredVariants, 5)
	for i, v := range filteredVariants {
		if i == 4 {
			assert.Equal(t, fmt.Sprintf("FilterWith()[%d].Format", i), v.Format, ImageFormatJPEG)
			assert.Equal(t, fmt.Sprintf("FilterWith()[%d].Width", i), v.Width, 320)
		} else {
			assert.Equal(t, fmt.Sprintf("FilterWith()[%d].Format", i), v.Format, ImageFormatWebP)
			assert.Equal(t, fmt.Sprintf("FilterWith()[%d].Width", i), v.Width, variants[i].Width)
		}
	}
}

func TestFilterWithKeepFallbackFalseDropsFallbackVariant(t *testing.T) {
	variants := Variants{
		makeVariant(1, ImageFormatWebP, 320),
		makeVariant(FallbackLayerID, ImageFormatJPEG, 320),
	}
	filteredVariants := variants.FilterWith(VariantFilterOptions{
		IncludeFormatSet: ImageFormatSet{
			ImageFormatWebP: {},
		},
		KeepFallback: false,
	})

	assert.LenIs(t, "FilterWith()", filteredVariants, 1)
	assert.Equal(t, "FilterWith()[0].Format", filteredVariants[0].Format, ImageFormatWebP)
	assert.Equal(t, "FilterWith()[0].Width", filteredVariants[0].Width, 320)
}
