package config

import "time"

type StaticFilesConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	RootDirectory string `mapstructure:"root_dir"`

	MaxAge               time.Duration `mapstructure:"max_age"`
	StaleWhileRevalidate time.Duration `mapstructure:"stale_while_revalidate"`
	Immutable            bool          `mapstructure:"immutable"`
	NoTransform          bool          `mapstructure:"no_transform"`
	NoSniff              bool          `mapstructure:"nosniff"`
}

type ServerConfig struct {
	Port uint16 `mapstructure:"port"`

	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes"`

	StaticFiles StaticFilesConfig `mapstructure:"static"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port: 4096,

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		ShutdownTimeout:   30 * time.Second,
		MaxHeaderBytes:    1 << 10, // 1 KiB

		StaticFiles: StaticFilesConfig{
			Enabled:              false,
			MaxAge:               3 * 24 * time.Hour,
			StaleWhileRevalidate: 2 * 24 * time.Hour,
			Immutable:            true,
			NoTransform:          false,
			NoSniff:              true,
		},
	}
}
