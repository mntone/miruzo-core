package media

type LayerIDType = uint32

var FallbackLayerID LayerIDType = 9

func IsVariantFallbackLayerID(layerID LayerIDType) bool {
	return layerID == FallbackLayerID
}

type QualityType = uint16

type Variant struct {
	RelativePath string       `json:"rel"`
	LayerID      LayerIDType  `json:"layer_id"`
	Format       string       `json:"format"`
	Codecs       string       `json:"codecs,omitempty"`
	Bytes        uint32       `json:"bytes"`
	Width        uint16       `json:"width"`
	Height       uint16       `json:"height"`
	Quality      *QualityType `json:"quality"`
}
