package config

import "github.com/mntone/miruzo-core/miruzo/internal/domain/media"

type APIConfig struct {
	VariantLayers media.VariantLayersSpec
	MediaPublic   MediaPublicConfig `mapstructure:"media_public"`
	Retry         RetryConfig       `mapstructure:"retry"`
}

func DefaultAPIConfig() APIConfig {
	return APIConfig{
		VariantLayers: DefaultVariantLayerConfig(),
		MediaPublic:   DefaultMediaPublicConfig(),
		Retry:         DefaultRetryConfig(),
	}
}
