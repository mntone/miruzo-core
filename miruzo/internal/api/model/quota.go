package model

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type Quota struct {
	Period    string         `json:"period"`
	ResetAt   time.Time      `json:"reset_at"`
	Limit     model.QuotaInt `json:"limit"`
	Remaining model.QuotaInt `json:"remaining"`
}
