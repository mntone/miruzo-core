package quota

import (
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/service/user"
)

var ErrInvalidPeriodType = errors.New("invalid period type")

func mapPeriodType(val user.PeriodType) (string, error) {
	switch val {
	case user.PeriodTypeDaily:
		return "daily", nil
	}

	return "", fmt.Errorf("%w: type=%d", ErrInvalidPeriodType, val)
}

func mapQuota(result user.QuotaResult) (quotaResponse, error) {
	lovePeriod, err := mapPeriodType(result.Love.Period)
	if err != nil {
		return quotaResponse{}, err
	}

	return quotaResponse{
		Love: quotaItem{
			Period:    lovePeriod,
			ResetAt:   result.Love.ResetAt,
			Limit:     result.Love.Limit,
			Remaining: result.Love.Remaining,
		},
	}, nil
}
