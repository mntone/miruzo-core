package config

type APIConfig struct {
	VariantLayers []VariantLayerConfig
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
