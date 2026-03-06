package app

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/spf13/viper"
)

func LoadConfig() (config.AppConfig, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return config.AppConfig{}, fmt.Errorf("read config: %w", err)
	}

	cfg := config.DefaultAppConfig()
	if err := viper.Unmarshal(&cfg); err != nil {
		return config.AppConfig{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
