package model

import (
	"time"

	"github.com/samber/mo"
)

type Job struct {
	Name       string
	StartedAt  time.Time
	FinishedAt mo.Option[time.Time]
}
