package user

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	repository          persist.UserRepository
	clk                 clock.Provider
	dailyPeriodResolver period.DailyResolver
	dailyLoveLimit      model.QuotaInt // <= model.MaxQuotaInt
}

func New(
	repo persist.UserRepository,
	clockProvider clock.Provider,
	dailyPeriodResolver period.DailyResolver,
	dailyLoveLimit model.QuotaInt,
) (Service, error) {
	if dailyLoveLimit < 1 || dailyLoveLimit > model.MaxQuotaInt {
		return Service{}, fmt.Errorf("invalid daily_love_limit: %d", dailyLoveLimit)
	}

	return Service{
		repository:          repo,
		clk:                 clockProvider,
		dailyPeriodResolver: dailyPeriodResolver,
		dailyLoveLimit:      dailyLoveLimit,
	}, nil
}
