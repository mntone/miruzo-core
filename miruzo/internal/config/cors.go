package config

import "time"

type CORSConfig struct {
	AllowOrigins []string       `mapstructure:"allow_origin"`
	MaxAge       *time.Duration `mapstructure:"max_age"`
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: nil,
		MaxAge:       nil,
	}
}
