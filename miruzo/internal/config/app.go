package config

type AppConfig struct {
	API      APIConfig
	Database DatabaseConfig `mapstructure:"database"`
	Period   PeriodConfig   `mapstructure:"period"`
	Quota    QuotaConfig    `mapstructure:"quota"`
	Server   ServerConfig   `mapstructure:"server"`
	Score    ScoreConfig    `mapstructure:"score"`
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		API:      DefaultAPIConfig(),
		Database: DefaultDatabaseConfig(),
		Period:   DefaultPeriodConfig(),
		Quota:    DefaultQuotaConfig(),
		Server:   DefaultServerConfig(),
		Score:    DefaultScoreConfig(),
	}
}
