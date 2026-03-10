package config

import (
	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type ScoreConfig struct {
	InitialScore          model.ScoreType `mapstructure:"initial_score"`
	MinimumScore          model.ScoreType `mapstructure:"minimum_score"`
	PublicMinimumScore    model.ScoreType `mapstructure:"public_minimum_score"`
	MaximumScore          model.ScoreType `mapstructure:"maximum_score"`
	EngagedScoreThreshold model.ScoreType `mapstructure:"engaged_score_threshold"`
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
