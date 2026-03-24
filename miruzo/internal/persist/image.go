package persist

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Image struct {
	IngestID   model.IngestIDType
	IngestedAt time.Time
	Type       model.ImageType

	Original Variant
	Fallback mo.Option[Variant]
	Variants []Variant
}

type ImageListCursor interface {
	~model.ScoreType | time.Time
}

type ImageWithCursor[C ImageListCursor] struct {
	Image  Image
	Cursor C
}

type ImageWithStats struct {
	Image
	Stats model.Stats
}
