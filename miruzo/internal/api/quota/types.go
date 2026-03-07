package quota

import "time"

type quotaItem struct {
	Period    string    `json:"period"`
	ResetAt   time.Time `json:"reset_at"`
	Limit     uint16    `json:"limit"`
	Remaining uint16    `json:"remaining"`
}

type quotaResponse struct {
	Love quotaItem `json:"love"`
}
