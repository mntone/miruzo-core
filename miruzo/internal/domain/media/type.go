package media

import "github.com/samber/mo"

type Variant struct {
	RelativePath string
	LayerID      LayerIDType
	Format       string
	Codecs       string
	Bytes        uint32
	Width        uint16
	Height       uint16
	Quality      mo.Option[QualityType]
}

func (v Variant) IsFallback() bool {
	return v.LayerID == FallbackLayerID
}
