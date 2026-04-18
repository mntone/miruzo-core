package model

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
)

type ImageType uint8

const (
	ImageTypeUnspecified ImageType = iota
	ImageTypePhoto
	ImageTypeIllust
	ImageTypeGraphic
	imageTypeLast
)

const (
	imageTypeStringPhoto   = "photo"
	imageTypeStringIllust  = "illust"
	imageTypeStringGraphic = "graphic"
)

func ValidateImageType[I constraints.Integer](value I) (ImageType, error) {
	if 0 > value || value >= I(imageTypeLast) {
		return 0, fmt.Errorf("%w: type=%d", errInvalidImageType, value)
	}

	return ImageType(value), nil
}

func (t ImageType) MarshalJSON() ([]byte, error) {
	switch t {
	case ImageTypePhoto:
		return json.Marshal(imageTypeStringPhoto)
	case ImageTypeIllust:
		return json.Marshal(imageTypeStringIllust)
	case ImageTypeGraphic:
		return json.Marshal(imageTypeStringGraphic)
	}
	return json.Marshal(nil)
}

func (t *ImageType) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		switch value {
		case imageTypeStringPhoto:
			*t = ImageTypePhoto
			return nil
		case imageTypeStringIllust:
			*t = ImageTypeIllust
			return nil
		case imageTypeStringGraphic:
			*t = ImageTypeGraphic
			return nil
		}
		return fmt.Errorf("%w: type=%s", errInvalidImageType, value)

	case nil:
		*t = ImageTypeUnspecified
		return nil
	}

	return fmt.Errorf("%w: type=%v", errInvalidImageType, v)
}

type Image struct {
	IngestID   IngestIDType
	IngestedAt time.Time
	Type       ImageType

	VariantBundle
}

type ImageListCursorScalar interface {
	ScoreType | time.Time
}

type ImageListCursorKey[ScalarType ImageListCursorScalar] struct {
	Primary   ScalarType
	Secondary IngestIDType
}

type ImageWithStats struct {
	Image
	Stats Stats
}
