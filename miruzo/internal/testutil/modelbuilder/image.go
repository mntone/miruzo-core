package modelbuilder

import (
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

type imageBuilder struct {
	BaseTime time.Time

	IngestID   model.IngestIDType
	IngestedAt time.Time
	ImageType  model.ImageType

	FilenameBase         string
	ImageAspectRatio     float64
	OriginalWidth        uint16
	PrimaryWidths        []uint16
	FallbackVariantWidth uint16
}

func Image(id model.IngestIDType) *imageBuilder {
	if id <= 0 {
		panic("invalid ingest id")
	}

	return &imageBuilder{
		BaseTime: defaultBaseTime,

		IngestID:   id,
		IngestedAt: defaultBaseTime,
		ImageType:  model.ImageTypeUnspecified,

		FilenameBase:         strconv.Itoa(int(id)),
		ImageAspectRatio:     4.0 / 3,
		OriginalWidth:        768,
		PrimaryWidths:        []uint16{320, 640},
		FallbackVariantWidth: 320,
	}
}

func (b *imageBuilder) Ingested(at time.Time) *imageBuilder {
	b.IngestedAt = at
	return b
}

func (b *imageBuilder) IngestedOffset(v any) *imageBuilder {
	if at, present := resolveOffsetTime(v, b.BaseTime).Get(); present {
		return b.Ingested(at)
	}
	return b
}

func (b *imageBuilder) Type(t model.ImageType) *imageBuilder {
	b.ImageType = t
	return b
}

func (b *imageBuilder) Filename(filename string) *imageBuilder {
	b.FilenameBase = filename
	return b
}

func (b *imageBuilder) AspectRatio(aspectRatio float64) *imageBuilder {
	if math.IsNaN(aspectRatio) || aspectRatio < 0.1 || aspectRatio > 10 {
		panic("invalid aspect ratio")
	}

	b.ImageAspectRatio = aspectRatio
	return b
}

func (b *imageBuilder) Original(width uint16) *imageBuilder {
	b.OriginalWidth = width
	return b
}

func (b *imageBuilder) AppendPrimaryVariant(width uint16) *imageBuilder {
	if width < 4 {
		panic("invalid width")
	}

	b.PrimaryWidths = append(b.PrimaryWidths, width)
	return b
}

func (b *imageBuilder) FallbackVariant(width uint16) *imageBuilder {
	b.FallbackVariantWidth = width
	return b
}

func (b *imageBuilder) Build() persist.Image {
	slices.Sort(b.PrimaryWidths)
	layers := make([]persist.Variant, len(b.PrimaryWidths)+1)
	for i, width := range b.PrimaryWidths {
		layers[i] = CreateVariant(
			1,
			media.ImageFormatWebP,
			b.FilenameBase,
			width,
			b.ImageAspectRatio,
		)
	}
	layers[len(layers)-1] = CreateVariant(
		media.FallbackLayerID,
		media.ImageFormatJPEG,
		b.FilenameBase,
		b.FallbackVariantWidth,
		b.ImageAspectRatio,
	)

	return persist.Image{
		IngestID:   b.IngestID,
		IngestedAt: b.IngestedAt,
		Type:       b.ImageType,
		Original: CreateVariant(
			0,
			media.ImageFormatPNG,
			b.FilenameBase,
			b.OriginalWidth,
			b.ImageAspectRatio,
		),
		Fallback: mo.None[persist.Variant](),
		Layers:   layers,
	}
}
