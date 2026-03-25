package media

import (
	"log"
	"slices"
	"sort"
)

type VariantLayersBuilder struct {
	spec VariantLayersSpec
}

func NewVariantLayerBuilder(spec VariantLayersSpec) *VariantLayersBuilder {
	return &VariantLayersBuilder{
		spec: spec,
	}
}

func (builder *VariantLayersBuilder) GroupVariantsByLayer(variants []Variant) VariantLayers {
	layered := make(map[LayerIDType][]Variant, len(builder.spec))

	for _, layer := range builder.spec {
		layered[layer.LayerID] = nil
	}

	unknownCount := 0
	unknownLayerIDs := make(map[LayerIDType]struct{})
	for _, variant := range variants {
		layer, ok := layered[variant.LayerID]
		if !ok {
			unknownCount++
			unknownLayerIDs[variant.LayerID] = struct{}{}
			continue
		}

		layered[variant.LayerID] = append(layer, variant)
	}
	if unknownCount > 0 {
		ids := make([]LayerIDType, 0, len(unknownLayerIDs))
		for layerID := range unknownLayerIDs {
			ids = append(ids, layerID)
		}
		slices.Sort(ids)
		log.Printf(
			"unknown variant layers dropped: count=%d layer_ids=%v",
			unknownCount,
			ids,
		)
	}

	for layerID := range layered {
		layer := layered[layerID]
		sort.Slice(layer, func(i, j int) bool {
			return layer[i].Width < layer[j].Width
		})
		layered[layerID] = layer
	}

	variantLayers := make(VariantLayers, 0, len(builder.spec))
	for _, layer := range builder.spec {
		if variants := layered[layer.LayerID]; variants != nil {
			variantLayers = append(variantLayers, variants)
		}
	}

	return variantLayers
}
