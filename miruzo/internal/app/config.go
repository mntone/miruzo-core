package app

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/spf13/viper"
)

//go:embed config.sample.yaml
var sampleConfigData []byte

func writeSampleConfigFile() error {
	f, err := os.Create("config.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(sampleConfigData)
	return err
}

func LoadConfig() (config.AppConfig, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, notFoundError := errors.AsType[viper.ConfigFileNotFoundError](err); notFoundError {
			if writeError := writeSampleConfigFile(); writeError != nil {
				return config.AppConfig{}, fmt.Errorf("write sample config: %w", writeError)
			}
		}

		return config.AppConfig{}, fmt.Errorf("read config: %w", err)
	}

	cfg := config.DefaultAppConfig()
	if err := viper.Unmarshal(&cfg); err != nil {
		return config.AppConfig{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return config.AppConfig{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}
