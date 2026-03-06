package config

import "time"

type RetryPolicyConfig struct {
	MaxAttempts uint32        `mapstructure:"max_attempts"`
	BaseDelay   time.Duration `mapstructure:"base_delay"`
	MinDelay    time.Duration `mapstructure:"min_delay"`
	MaxDelay    time.Duration `mapstructure:"max_delay"`
}

type RetryConfig struct {
	Read  RetryPolicyConfig `mapstructure:"read"`
	Write RetryPolicyConfig `mapstructure:"write"`
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		Read: RetryPolicyConfig{
			MaxAttempts: 3,
			BaseDelay:   20 * time.Millisecond,
			MinDelay:    10 * time.Millisecond,
			MaxDelay:    333 * time.Millisecond,
		},
		Write: RetryPolicyConfig{
			MaxAttempts: 1,
		},
	}
}
