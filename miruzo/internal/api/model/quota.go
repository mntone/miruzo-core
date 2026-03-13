package model

import "time"

type Quota struct {
	Period    string    `json:"period"`
	ResetAt   time.Time `json:"reset_at"`
	Limit     int16     `json:"limit"`
	Remaining int16     `json:"remaining"`
}
