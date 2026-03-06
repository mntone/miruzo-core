package app

import (
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

func newBackoffPolicyFromConfig(cfg config.RetryPolicyConfig) backoff.Policy {
	return &backoff.ExponentialPolicy{
		MaxAttempts:  cfg.MaxAttempts,
		BaseDelay:    cfg.BaseDelay,
		MinimumDelay: cfg.MinDelay,
		MaximumDelay: cfg.MaxDelay,
		Jitter:       backoff.JitterEqual,
	}
}
