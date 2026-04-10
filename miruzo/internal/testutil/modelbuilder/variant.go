package modelbuilder

import (
	"fmt"
	"math"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

func CreateVariant(
	layerID media.LayerIDType,
	format media.ImageFormat,
	filename string,
	width uint16,
	aspectRatio float64,
) persist.Variant {
	var dirname string
	if layerID == 0 {
		dirname = "l0orig"
	} else {
		dirname = fmt.Sprintf("l%dw%d", layerID, width)
	}

	var codecs string
	if format == media.ImageFormatWebP {
		codecs = "vp8"
	}

	return persist.Variant{
		RelativePath: fmt.Sprintf("%s/%s.%s", dirname, filename, format.String()),
		LayerID:      layerID,
		Format:       format,
		Codecs:       codecs,
		Bytes:        234,
		Width:        width,
		Height:       uint16((int32(math.Ceil(float64(width)/aspectRatio)) + 1) &^ 1),
		Quality:      nil,
	}
}
