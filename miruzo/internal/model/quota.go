package model

import "time"

type PeriodType uint8

const (
	PeriodTypeUnspecified PeriodType = iota
	PeriodTypeDaily
)

type Quota struct {
	Period    PeriodType
	ResetAt   time.Time
	Limit     int16
	Remaining int16
}
