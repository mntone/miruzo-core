package config

type QuotaConfig struct {
	DailyLoveLimit int16 `mapstructure:"daily_love_limit"`
}

func DefaultQuotaConfig() QuotaConfig {
	return QuotaConfig{
		DailyLoveLimit: 3,
	}
}
