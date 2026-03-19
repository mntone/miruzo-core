package reaction

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type statModel struct {
	Score        model.ScoreType `json:"score"`
	FirstLovedAt *time.Time      `json:"first_loved_at"`
	LastLovedAt  *time.Time      `json:"last_loved_at"`
}

type loveResponse struct {
	Quota model.Quota `json:"quota"`
	Stats statModel   `json:"stats"`
}
