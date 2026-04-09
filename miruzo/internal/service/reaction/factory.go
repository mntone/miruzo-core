package reaction

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	prov                     persist.PersistenceProvider
	clk                      clock.Provider
	dailyPeriodResolver      period.DailyResolver
	scoreCalculator          score.Calculator
	dailyLoveLimit           model.QuotaInt // <= model.MaxQuotaInt
	hallOfFameScoreThreshold model.ScoreType
}

func New(
	persistenceProvider persist.PersistenceProvider,
	clockProvider clock.Provider,
	dailyPeriodResolver period.DailyResolver,
	scoreCalculator score.Calculator,
	dailyLoveLimit model.QuotaInt,
	hallOfFameScoreThreshold model.ScoreType,
) (*Service, error) {
	if dailyLoveLimit < 1 || dailyLoveLimit > model.MaxQuotaInt {
		return nil, fmt.Errorf("invalid daily_love_limit: %d", dailyLoveLimit)
	}

	return &Service{
		prov:                     persistenceProvider,
		clk:                      clockProvider,
		dailyPeriodResolver:      dailyPeriodResolver,
		scoreCalculator:          scoreCalculator,
		dailyLoveLimit:           dailyLoveLimit,
		hallOfFameScoreThreshold: hallOfFameScoreThreshold,
	}, nil
}
