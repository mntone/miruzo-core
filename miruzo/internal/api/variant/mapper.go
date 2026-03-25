package variant

import (
	"math"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

func toManbytes(sizeInBytes uint32) uint16 {
	return uint16(math.Ceil(float64(sizeInBytes) / 10_000))
}

func mapVariant(
	variant media.Variant,
	mediaURLBuilder MediaURLBuilder,
) VariantModel {
	return VariantModel{
		Source:   mediaURLBuilder.Build(variant.RelativePath),
		Format:   variant.Format,
		Codecs:   variant.Codecs,
		Manbytes: toManbytes(variant.Bytes),
		Width:    variant.Width,
		Height:   variant.Height,
	}
}

func mapNullableVariant(
	optionalEntry mo.Option[media.Variant],
	mediaURLBuilder MediaURLBuilder,
) *VariantModel {
	entry, present := optionalEntry.Get()
	if !present {
		return nil
	}

	model := mapVariant(entry, mediaURLBuilder)
	return &model
}

func mapVariantsToLayers(
	variants [][]media.Variant,
	mediaURLBuilder MediaURLBuilder,
) [][]VariantModel {
	result := make([][]VariantModel, len(variants))
	for i, layer := range variants {
		result[i] = make([]VariantModel, len(layer))
		for j, variant := range layer {
			result[i][j] = mapVariant(variant, mediaURLBuilder)
		}
	}

	return result
}

func MapVariantLayers(
	entry model.Image,
	mediaURLBuilder MediaURLBuilder,
) VariantLayersModel {
	return VariantLayersModel{
		Original: mapVariant(entry.Original, mediaURLBuilder),
		Fallback: mapNullableVariant(entry.Fallback, mediaURLBuilder),
		Layers:   mapVariantsToLayers(entry.Layers, mediaURLBuilder),
	}
}
