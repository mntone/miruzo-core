package user

import (
	"context"
	"time"
)

type PeriodType uint8

const (
	PeriodTypeUnspecified PeriodType = iota
	PeriodTypeDaily
)

type quotaItem struct {
	Period    PeriodType
	ResetAt   time.Time
	Limit     uint16
	Remaining uint16
}

type QuotaResult struct {
	Love quotaItem
}

func (srv Service) GetQuota(
	requestContext context.Context,
) (QuotaResult, error) {
	user, err := srv.repository.GetSingletonUser(requestContext)
	if err != nil {
		return QuotaResult{}, err
	}

	loveUsed := uint16(user.DailyLoveUsed)
	loveRemaining := uint16(0)
	if loveUsed < srv.dailyLoveLimit {
		loveRemaining = srv.dailyLoveLimit - loveUsed
	}

	result := QuotaResult{
		Love: quotaItem{
			Period:    PeriodTypeDaily,
			ResetAt:   srv.dailyPeriodResolver.PeriodEnd(time.Now()).UTC(),
			Limit:     srv.dailyLoveLimit,
			Remaining: loveRemaining,
		},
	}
	return result, nil
}
