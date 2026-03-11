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

	// --- view ---

	ViewBonusAtFirst  model.ScoreType            `mapstructure:"view_bonus_at_first"`
	ViewBonusByDays   []model.ScoreViewBonusRule `mapstructure:"view_bonus_by_days"`
	ViewBonusFallback model.ScoreType            `mapstructure:"view_bonus_fallback"`

	// --- memo ---

	MemoBonus   model.ScoreType `mapstructure:"memo_bonus"`
	MemoPenalty model.ScoreType `mapstructure:"memo_penalty"`

	// --- love ---

	LoveBonus   model.ScoreType `mapstructure:"love_bonus"`
	LovePenalty model.ScoreType `mapstructure:"love_penalty"`
}

func DefaultScoreConfig() ScoreConfig {
	return ScoreConfig{
		InitialScore:          100,
		MinimumScore:          -20000,
		PublicMinimumScore:    0,
		MaximumScore:          200,
		EngagedScoreThreshold: 160,

		// --- view ---

		ViewBonusAtFirst: 10,
		ViewBonusByDays: []model.ScoreViewBonusRule{
			{Days: 1, Bonus: 10},
			{Days: 3, Bonus: 7},
			{Days: 7, Bonus: 5},
			{Days: 30, Bonus: 3},
		},
		ViewBonusFallback: 2,

		// --- memo ---

		MemoBonus:   1,
		MemoPenalty: -1,

		// --- love ---

		LoveBonus:   20,
		LovePenalty: -18,
	}
}
