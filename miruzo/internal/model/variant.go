package model

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/samber/mo"
)

type VariantBundle struct {
	Original media.Variant
	Fallback mo.Option[media.Variant]
	Layers   media.VariantLayers
}
