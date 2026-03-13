package reaction

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	mgr                 persist.PersistenceManager
	dailyPeriodResolver period.DailyResolver
	scoreCalculator     score.Calculator
	dailyLoveLimit      int16
}

func New(
	persistenceManager persist.PersistenceManager,
	dailyPeriodResolver period.DailyResolver,
	scoreCalculator score.Calculator,
	dailyLoveLimit int16,
) Service {
	return Service{
		mgr:                 persistenceManager,
		dailyPeriodResolver: dailyPeriodResolver,
		scoreCalculator:     scoreCalculator,
		dailyLoveLimit:      dailyLoveLimit,
	}
}
