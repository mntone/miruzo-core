package score

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/samber/mo"
)

type Calculator struct {
	dailyResolver period.DailyResolver

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
	viewBonusAtFirst model.ScoreType,
	viewBonusByDays []model.ScoreViewBonusRule,
	viewBonusFallback model.ScoreType,
	memoBonus model.ScoreType,
	memoPenalty model.ScoreType,
	loveBonus model.ScoreType,
	lovePenalty model.ScoreType,
) Calculator {
	return Calculator{
		dailyResolver:     dailyResolver,
		viewBonusAtFirst:  viewBonusAtFirst,
		viewBonusByDays:   viewBonusByDays,
		viewBonusFallback: viewBonusFallback,
		memoBonus:         memoBonus,
		memoPenalty:       memoPenalty,
		loveBonus:         loveBonus,
		lovePenalty:       lovePenalty,
	}
}

func (calc Calculator) ViewDelta(
	lastViewedAt mo.Option[time.Time],
	evaluatedAt time.Time,
) (delta model.ScoreType, negative bool) {
	validLastView, present := lastViewedAt.Get()
	if !present {
		return calc.viewBonusAtFirst, false
	}

	evaluate := calc.dailyResolver.PeriodStart(evaluatedAt)
	lastView := calc.dailyResolver.PeriodStart(validLastView)

	days := int32(evaluate.Sub(lastView) / (24 * time.Hour))
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
