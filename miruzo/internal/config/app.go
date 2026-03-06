package config

type AppConfig struct {
	API      APIConfig
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	Score    ScoreConfig    `mapstructure:"score"`
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		API:      DefaultAPIConfig(),
		Database: DefaultDatabaseConfig(),
		Server:   DefaultServerConfig(),
		Score:    DefaultScoreConfig(),
	}
}
