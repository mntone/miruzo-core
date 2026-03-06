package config

type ScoreConfig struct {
	InitialScore          int16 `mapstructure:"initial_score"`
	MinimumScore          int16 `mapstructure:"minimum_score"`
	PublicMinimumScore    int16 `mapstructure:"public_minimum_score"`
	MaximumScore          int16 `mapstructure:"maximum_score"`
	EngagedScoreThreshold int16 `mapstructure:"engaged_score_threshold"`
}

func DefaultScoreConfig() ScoreConfig {
	return ScoreConfig{
		InitialScore:          100,
		MinimumScore:          -20000,
		PublicMinimumScore:    0,
		MaximumScore:          200,
		EngagedScoreThreshold: 160,
	}
}
