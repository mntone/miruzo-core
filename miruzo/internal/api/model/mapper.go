package model

import (
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

var ErrInvalidPeriodType = errors.New("invalid period type")

func mapPeriodType(value model.PeriodType) (string, error) {
	switch value {
	case model.PeriodTypeDaily:
		return "daily", nil
	}

	return "", fmt.Errorf("%w: type=%d", ErrInvalidPeriodType, value)
}

func MapQuota(q model.Quota) (Quota, error) {
	period, err := mapPeriodType(q.Period)
	if err != nil {
		return Quota{}, err
	}

	return Quota{
		Period:    period,
		ResetAt:   q.ResetAt,
		Limit:     q.Limit,
		Remaining: q.Remaining,
	}, nil
}
