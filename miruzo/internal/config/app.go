package config

import "fmt"

type AppConfig struct {
	API      APIConfig
	Database DatabaseConfig `mapstructure:"database"`
	Period   PeriodConfig   `mapstructure:"period"`
	Quota    QuotaConfig    `mapstructure:"quota"`
	Server   ServerConfig   `mapstructure:"server"`
	Score    ScoreConfig    `mapstructure:"score"`
	View     ViewConfig     `mapstructure:"view"`
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		API:      DefaultAPIConfig(),
		Database: DefaultDatabaseConfig(),
		Period:   DefaultPeriodConfig(),
		Quota:    DefaultQuotaConfig(),
		Server:   DefaultServerConfig(),
		Score:    DefaultScoreConfig(),
		View:     DefaultViewConfig(),
	}
}

func (c *AppConfig) Validate() error {
	err := c.Quota.Validate()
	if err != nil {
		return fmt.Errorf("invalid config: quota.%w", err)
	}

	return err
}
