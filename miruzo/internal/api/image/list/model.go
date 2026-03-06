package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type ImageListModel struct {
	IngestID persist.IngestID `json:"id"`
	variant.VariantLayersModel
}
