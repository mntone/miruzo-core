package domain

import (
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
)

func NewTestScoreCalculator(dailyResolver period.DailyResolver) score.Calculator {
	cfg := config.DefaultScoreConfig()
	return score.New(
		dailyResolver,
		cfg.DailyDecayNoAccessAdjustment,
		cfg.DailyDecayPenalty,
		cfg.DailyDecayInterval10dPenalty,
		cfg.DailyDecayHighScorePenalty,
		cfg.DailyDecayHighScoreThreshold,
		cfg.ViewBonusAtFirst,
		cfg.ViewBonusByDays,
		cfg.ViewBonusFallback,
		cfg.MemoBonus,
		cfg.MemoPenalty,
		cfg.LoveBonus,
		cfg.LovePenalty,
	)
}
