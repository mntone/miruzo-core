package shared

import (
	"encoding/json"

	"github.com/mntone/miruzo-core/miruzo/internal/model/media"
	"github.com/samber/mo"
)

func MapVariant(raw []byte) (media.Variant, error) {
	var variant media.Variant
	if err := json.Unmarshal(raw, &variant); err != nil {
		return media.Variant{}, err
	}

	return variant, nil
}

func MapNullableVariant(raw *[]byte) (mo.Option[media.Variant], error) {
	if raw == nil {
		return mo.None[media.Variant](), nil
	}

	var variant media.Variant
	if err := json.Unmarshal(*raw, &variant); err != nil {
		return mo.None[media.Variant](), err
	}

	return mo.Some(variant), nil
}

func MapVariants(raw []byte) ([]media.Variant, error) {
	var variants []media.Variant
	if err := json.Unmarshal(raw, &variants); err != nil {
		return nil, err
	}

	return variants, nil
}
