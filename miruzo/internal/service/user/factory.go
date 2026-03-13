package user

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	repository          persist.UserRepository
	dailyPeriodResolver period.DailyResolver
	dailyLoveLimit      int16
}

func New(
	repo persist.UserRepository,
	dailyPeriodResolver period.DailyResolver,
	dailyLoveLimit int16,
) Service {
	return Service{
		repository:          repo,
		dailyPeriodResolver: dailyPeriodResolver,
		dailyLoveLimit:      dailyLoveLimit,
	}
}
