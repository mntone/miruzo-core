package model

import "time"

type PeriodType uint8

const (
	PeriodTypeUnspecified PeriodType = iota
	PeriodTypeDaily
)

type QuotaInt int32

// Keep in sync with database CHECK constraints (users.daily_love_used).
const MaxQuotaInt QuotaInt = 100

type Quota struct {
	Period    PeriodType
	ResetAt   time.Time
	Limit     QuotaInt
	Remaining QuotaInt
}
