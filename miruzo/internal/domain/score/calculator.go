package score

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Calculator struct {
	dailyResolver period.DailyResolver

	// --- daily decay ---

	dailyDecayNoAccessAdjustment model.ScoreType
	dailyDecayPenalty            model.ScoreType
	dailyDecayInterval10dPenalty model.ScoreType
	dailyDecayHighScorePenalty   model.ScoreType
	dailyDecayHighScoreThreshold model.ScoreType

	// --- view ---

	viewBonusAtFirst  model.ScoreType
	viewBonusByDays   []model.ScoreViewBonusRule
	viewBonusFallback model.ScoreType

	// --- memo ---

	memoBonus   model.ScoreType
	memoPenalty model.ScoreType

	// --- love ---

	loveBonus   model.ScoreType
	lovePenalty model.ScoreType
}

func New(
	dailyResolver period.DailyResolver,
	dailyDecayNoAccessAdjustment model.ScoreType,
	dailyDecayPenalty model.ScoreType,
	dailyDecayInterval10dPenalty model.ScoreType,
	dailyDecayHighScorePenalty model.ScoreType,
	dailyDecayHighScoreThreshold model.ScoreType,
	viewBonusAtFirst model.ScoreType,
	viewBonusByDays []model.ScoreViewBonusRule,
	viewBonusFallback model.ScoreType,
	memoBonus model.ScoreType,
	memoPenalty model.ScoreType,
	loveBonus model.ScoreType,
	lovePenalty model.ScoreType,
) Calculator {
	return Calculator{
		dailyResolver:                dailyResolver,
		dailyDecayNoAccessAdjustment: dailyDecayNoAccessAdjustment,
		dailyDecayPenalty:            dailyDecayPenalty,
		dailyDecayInterval10dPenalty: dailyDecayInterval10dPenalty,
		dailyDecayHighScorePenalty:   dailyDecayHighScorePenalty,
		dailyDecayHighScoreThreshold: dailyDecayHighScoreThreshold,
		viewBonusAtFirst:             viewBonusAtFirst,
		viewBonusByDays:              viewBonusByDays,
		viewBonusFallback:            viewBonusFallback,
		memoBonus:                    memoBonus,
		memoPenalty:                  memoPenalty,
		loveBonus:                    loveBonus,
		lovePenalty:                  lovePenalty,
	}
}

func (calc Calculator) calcDays(
	lastViewedAt time.Time,
	evaluatedAt time.Time,
) int32 {
	evaluate := calc.dailyResolver.PeriodStart(evaluatedAt)
	lastView := calc.dailyResolver.PeriodStart(lastViewedAt)
	return int32(evaluate.Sub(lastView) / (24 * time.Hour))
}

func (calc Calculator) DailyDecay(
	score model.ScoreType,
	lastViewedAt time.Time,
	evaluatedAt time.Time,
) (newScore model.ScoreType, negative bool) {
	newScore = score

	days := calc.calcDays(lastViewedAt, evaluatedAt)
	negative = days < 0

	if score >= calc.dailyDecayHighScoreThreshold {
		newScore += calc.dailyDecayHighScorePenalty
	} else if days != 0 && days%10 == 0 {
		newScore += calc.dailyDecayInterval10dPenalty
	} else {
		newScore += calc.dailyDecayPenalty
	}

	if days > 1 {
		newScore += calc.dailyDecayNoAccessAdjustment
	}

	return
}

func (calc Calculator) ViewDelta(
	lastViewedAt mo.Option[time.Time],
	evaluatedAt time.Time,
) (delta model.ScoreType, negative bool) {
	validLastView, present := lastViewedAt.Get()
	if !present {
		return calc.viewBonusAtFirst, false
	}

	days := calc.calcDays(validLastView, evaluatedAt)
	if days > 0 {
		for _, rule := range calc.viewBonusByDays {
			if days <= rule.Days {
				return rule.Bonus, false
			}
		}
		return calc.viewBonusFallback, false
	}

	return 0, days < 0
}

func (calc Calculator) MemoAddedDelta() model.ScoreType {
	return calc.memoBonus
}

func (calc Calculator) MemoRemovedDelta() model.ScoreType {
	return calc.memoPenalty
}

func (calc Calculator) LoveDelta() model.ScoreType {
	return calc.loveBonus
}

func (calc Calculator) LoveCanceledDelta() model.ScoreType {
	return calc.lovePenalty
}
