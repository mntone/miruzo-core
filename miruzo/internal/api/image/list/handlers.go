package list

import (
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

type handler struct {
	service             imagelist.Service
	variantLayersConfig []config.VariantLayerConfig
	mediaURLBuilder     variant.MediaURLBuilder
}

func NewHandler(
	srv imagelist.Service,
	variantLayersConfig []config.VariantLayerConfig,
	mediaURLBuilder variant.MediaURLBuilder,
) *handler {
	return &handler{
		service:             srv,
		variantLayersConfig: variantLayersConfig,
		mediaURLBuilder:     mediaURLBuilder,
	}
}
