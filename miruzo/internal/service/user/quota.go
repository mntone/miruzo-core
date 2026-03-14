package user

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

type QuotaResult struct {
	Love model.Quota
}

func (srv Service) GetQuota(
	requestContext context.Context,
) (QuotaResult, error) {
	current := srv.clk.Now()

	user, err := srv.repository.GetSingletonUser(requestContext)
	if err != nil {
		return QuotaResult{}, serviceerror.MapPersistError(err)
	}

	loveUsed := user.DailyLoveUsed
	loveRemaining := int16(0)
	if loveUsed < srv.dailyLoveLimit {
		loveRemaining = srv.dailyLoveLimit - loveUsed
	}

	result := QuotaResult{
		Love: model.Quota{
			Period:    model.PeriodTypeDaily,
			ResetAt:   srv.dailyPeriodResolver.PeriodEnd(current),
			Limit:     srv.dailyLoveLimit,
			Remaining: loveRemaining,
		},
	}
	return result, nil
}
