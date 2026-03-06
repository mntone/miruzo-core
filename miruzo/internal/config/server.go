package config

import "time"

type ServerConfig struct {
	Port uint16 `mapstructure:"port"`

	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:              4096,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		ShutdownTimeout:   30 * time.Second,
		MaxHeaderBytes:    1 << 12, // 4 KiB
	}
}
