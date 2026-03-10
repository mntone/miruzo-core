package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type ImageListModel struct {
	IngestID model.IngestIDType `json:"id"`
	variant.VariantLayersModel
}
