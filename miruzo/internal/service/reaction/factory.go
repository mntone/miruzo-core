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
	mgr                 persist.PersistenceManager
	clk                 clock.Provider
	dailyPeriodResolver period.DailyResolver
	scoreCalculator     score.Calculator
	dailyLoveLimit      model.QuotaInt // <= model.MaxQuotaInt
}

func New(
	persistenceManager persist.PersistenceManager,
	clockProvider clock.Provider,
	dailyPeriodResolver period.DailyResolver,
	scoreCalculator score.Calculator,
	dailyLoveLimit model.QuotaInt,
) (Service, error) {
	if dailyLoveLimit < 1 || dailyLoveLimit > model.MaxQuotaInt {
		return Service{}, fmt.Errorf("invalid daily_love_limit: %d", dailyLoveLimit)
	}

	return Service{
		mgr:                 persistenceManager,
		clk:                 clockProvider,
		dailyPeriodResolver: dailyPeriodResolver,
		scoreCalculator:     scoreCalculator,
		dailyLoveLimit:      dailyLoveLimit,
	}, nil
}
