package variant

import (
	"math"
	"sort"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

func toManbytes(sizeInBytes uint32) uint16 {
	return uint16(math.Ceil(float64(sizeInBytes) / 10_000))
}

func mapVariant(
	variant persist.Variant,
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
	optionalEntry mo.Option[persist.Variant],
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
	variants []persist.Variant,
	spec media.VariantLayersSpec,
	mediaURLBuilder MediaURLBuilder,
) [][]VariantModel {
	layered := make(map[media.LayerIDType][]VariantModel, len(spec))

	for _, layer := range spec {
		layered[layer.LayerID] = []VariantModel{}
	}

	for _, variantEntry := range variants {
		layer, ok := layered[variantEntry.LayerID]
		if !ok {
			continue
		}

		variantModel := mapVariant(variantEntry, mediaURLBuilder)

		layered[variantEntry.LayerID] = append(layer, variantModel)
	}

	for layerID := range layered {
		layer := layered[layerID]
		sort.Slice(layer, func(i, j int) bool {
			return layer[i].Width < layer[j].Width
		})
		layered[layerID] = layer
	}

	result := make([][]VariantModel, 0, len(spec))
	for _, layer := range spec {
		if entries := layered[layer.LayerID]; len(entries) > 0 {
			result = append(result, entries)
		}
	}

	return result
}

func MapVariantLayers(
	entry persist.Image,
	spec media.VariantLayersSpec,
	mediaURLBuilder MediaURLBuilder,
) VariantLayersModel {
	return VariantLayersModel{
		Original: mapVariant(entry.Original, mediaURLBuilder),
		Fallback: mapNullableVariant(entry.Fallback, mediaURLBuilder),
		Variants: mapVariantsToLayers(entry.Variants, spec, mediaURLBuilder),
	}
}
