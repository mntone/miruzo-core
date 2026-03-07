package config

import "time"

type PeriodConfig struct {
	DayStartOffset time.Duration `mapstructure:"day_start_offset"`
}

func DefaultPeriodConfig() PeriodConfig {
	return PeriodConfig{
		DayStartOffset: 5 * time.Hour,
	}
}
