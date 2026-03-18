package config

import (
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type QuotaConfig struct {
	DailyLoveLimit model.QuotaInt `mapstructure:"daily_love_limit"`
}

func DefaultQuotaConfig() QuotaConfig {
	return QuotaConfig{
		DailyLoveLimit: 3,
	}
}

func (c *QuotaConfig) Validate() error {
	if c.DailyLoveLimit < 1 || c.DailyLoveLimit > model.MaxQuotaInt {
		return errors.New("daily_love_limit")
	}
	return nil
}
