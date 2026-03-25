package model

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/samber/mo"
)

type VariantBundle struct {
	// Original is the canonical source variant.
	Original media.Variant
	// Fallback is an optional compatibility variant.
	Fallback mo.Option[media.Variant]
	// Layers groups variants by layer and size.
	Layers media.VariantLayers
}
