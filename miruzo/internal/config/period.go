package config

import "time"

type PeriodConfig struct {
	InitialLocation *string        `mapstructure:"initial_location"`
	DayStartOffset  *time.Duration `mapstructure:"day_start_offset"`
}

func DefaultPeriodConfig() PeriodConfig {
	return PeriodConfig{
		InitialLocation: nil,
		DayStartOffset:  nil,
	}
}
