package user

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	repository          persist.UserRepository
	clk                 clock.Provider
	dailyPeriodResolver period.DailyResolver
	dailyLoveLimit      int16
}

func New(
	repo persist.UserRepository,
	clockProvider clock.Provider,
	dailyPeriodResolver period.DailyResolver,
	dailyLoveLimit int16,
) Service {
	return Service{
		repository:          repo,
		clk:                 clockProvider,
		dailyPeriodResolver: dailyPeriodResolver,
		dailyLoveLimit:      dailyLoveLimit,
	}
}
