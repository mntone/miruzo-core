package shared

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeIntervalFromDuration(value time.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: value.Microseconds(),
		Valid:        true,
	}
}
