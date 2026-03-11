package config

type ViewConfig struct {
	Milestones []int64 `mapstructure:"milestones"`
}

func DefaultViewConfig() ViewConfig {
	return ViewConfig{
		Milestones: []int64{
			100,
			1_000,
			10_000,
			100_000,
			1_000_000,
			10_000_000,
			100_000_000,
			1_000_000_000,
			10_000_000_000,
			100_000_000_000,
			1_000_000_000_000,
		},
	}
}
