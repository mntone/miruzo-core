package model

import (
	"encoding/json"
	"errors"
	"time"
)

type PeriodType uint8

const (
	PeriodTypeUnspecified PeriodType = iota
	PeriodTypeDaily
)

func (p PeriodType) MarshalJSON() ([]byte, error) {
	switch p {
	case PeriodTypeDaily:
		return []byte("\"daily\""), nil
	case PeriodTypeUnspecified:
		return []byte("null"), nil
	default:
		return nil, errors.New("invalid period")
	}
}

func (p *PeriodType) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*p = PeriodTypeUnspecified
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "daily":
		*p = PeriodTypeDaily
	default:
		return errors.New("unknown period")
	}

	return nil
}

type QuotaInt int32

// Keep in sync with database CHECK constraints (users.daily_love_used).
const MaxQuotaInt QuotaInt = 100

type Quota struct {
	Period    PeriodType `json:"period,omitempty"`
	ResetAt   time.Time  `json:"reset_at"`
	Limit     QuotaInt   `json:"limit"`
	Remaining QuotaInt   `json:"remaining"`
}
