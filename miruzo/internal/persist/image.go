package persist

import (
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/samber/mo"
	"golang.org/x/exp/constraints"
)

type ImageType uint8

const (
	ImageTypeUnspecified ImageType = iota
	ImageTypePhoto
	ImageTypeIllust
	ImageTypeGraphic
)

func ParseImageType[T constraints.Integer](value T) (ImageType, error) {
	imageType := ImageType(value)

	switch imageType {
	case ImageTypeUnspecified, ImageTypePhoto, ImageTypeIllust, ImageTypeGraphic:
		return imageType, nil

	default:
		return 0, fmt.Errorf("%w: kind=%d", ErrInvalidImageKind, value)
	}
}

type Image struct {
	IngestID   IngestID
	IngestedAt time.Time
	Type       ImageType

	Original media.Variant
	Fallback mo.Option[media.Variant]
	Variants []media.Variant
}

type ImageListCursor interface {
	~int16 | time.Time
}

type ImageWithCursor[C ImageListCursor] struct {
	Image  Image
	Cursor C
}
