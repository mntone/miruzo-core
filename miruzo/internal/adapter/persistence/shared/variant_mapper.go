package shared

import (
	"encoding/json"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

func MapVariant(raw []byte) (persist.Variant, error) {
	var variant persist.Variant
	if err := json.Unmarshal(raw, &variant); err != nil {
		return persist.Variant{}, err
	}

	return variant, nil
}

func MapNullableVariant(raw *[]byte) (mo.Option[persist.Variant], error) {
	if raw == nil {
		return mo.None[persist.Variant](), nil
	}

	var variant persist.Variant
	if err := json.Unmarshal(*raw, &variant); err != nil {
		return mo.None[persist.Variant](), err
	}

	return mo.Some(variant), nil
}

func MapVariants(raw []byte) ([]persist.Variant, error) {
	var variants []persist.Variant
	if err := json.Unmarshal(raw, &variants); err != nil {
		return nil, err
	}

	return variants, nil
}
