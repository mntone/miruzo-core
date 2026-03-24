package persist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
)

type Variant struct {
	RelativePath string             `json:"rel"`
	LayerID      media.LayerIDType  `json:"layer_id"`
	Format       string             `json:"format"`
	Codecs       string             `json:"codecs,omitempty"`
	Bytes        uint32             `json:"bytes"`
	Width        uint16             `json:"width"`
	Height       uint16             `json:"height"`
	Quality      *media.QualityType `json:"quality"`
}
