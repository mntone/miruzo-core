package config

type MediaPublicConfig struct {
	BasePath string `mapstructure:"base_path"`
}

func DefaultMediaPublicConfig() MediaPublicConfig {
	return MediaPublicConfig{
		BasePath: "/media/",
	}
}
