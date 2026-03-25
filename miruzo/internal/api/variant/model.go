package variant

import "github.com/mntone/miruzo-core/miruzo/internal/domain/media"

type VariantModel struct {
	Source   string            `json:"src"`
	Format   media.ImageFormat `json:"format"`
	Codecs   string            `json:"codecs,omitempty"`
	Manbytes uint16            `json:"manbytes"`
	Width    uint16            `json:"w"`
	Height   uint16            `json:"h"`
}

type VariantLayersModel struct {
	Original VariantModel     `json:"original"`
	Fallback *VariantModel    `json:"fallback,omitempty"`
	Layers   [][]VariantModel `json:"variants"`
}
