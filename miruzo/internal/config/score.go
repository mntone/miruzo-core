package config

import (
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type ScoreConfig struct {
	InitialScore          model.ScoreType `mapstructure:"initial_score"`
	MinimumScore          model.ScoreType `mapstructure:"minimum_score"`
	PublicMinimumScore    model.ScoreType `mapstructure:"public_minimum_score"`
	MaximumScore          model.ScoreType `mapstructure:"maximum_score"`
	EngagedScoreThreshold model.ScoreType `mapstructure:"engaged_threshold"`

	// --- daily decay ---

	DailyDecayNoAccessAdjustment model.ScoreType `mapstructure:"daily_decay_no_access_adjustment"`
	DailyDecayPenalty            model.ScoreType `mapstructure:"daily_decay_penalty"`
	DailyDecayInterval10dPenalty model.ScoreType `mapstructure:"daily_decay_interval10d_penalty"`
	DailyDecayHighScorePenalty   model.ScoreType `mapstructure:"daily_decay_high_score_penalty"`
	DailyDecayHighScoreThreshold model.ScoreType `mapstructure:"daily_decay_high_score_threshold"`

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

	// --- hall of fame ---

	HallOfFameScoreThreshold model.ScoreType `mapstructure:"hof_threshold"`
}

func DefaultScoreConfig() ScoreConfig {
	return ScoreConfig{
		InitialScore:          100,
		MinimumScore:          -20000,
		PublicMinimumScore:    0,
		MaximumScore:          200,
		EngagedScoreThreshold: 160,

		// --- daily decay ---

		DailyDecayNoAccessAdjustment: 1,
		DailyDecayPenalty:            -2,
		DailyDecayInterval10dPenalty: -3,
		DailyDecayHighScorePenalty:   -3,
		DailyDecayHighScoreThreshold: 180, // == HallOfFameScoreThreshold

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

		// --- hall of fame ---

		HallOfFameScoreThreshold: 180,
	}
}

func (c *ScoreConfig) Validate() error {
	if c.MinimumScore > c.PublicMinimumScore ||
		c.PublicMinimumScore > c.MaximumScore {
		return errors.New("public_minimum_score")
	}
	if c.PublicMinimumScore > c.InitialScore ||
		c.InitialScore > c.MaximumScore {
		return errors.New("initial_score")
	}
	if c.PublicMinimumScore > c.EngagedScoreThreshold ||
		c.EngagedScoreThreshold > c.MaximumScore {
		return errors.New("engaged_threshold")
	}
	if c.PublicMinimumScore > c.HallOfFameScoreThreshold ||
		c.HallOfFameScoreThreshold > c.MaximumScore {
		return errors.New("hof_threshold")
	}
	return nil
}
