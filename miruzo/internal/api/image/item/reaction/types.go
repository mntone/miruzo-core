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

// Response payload for love actions.
type loveResponse struct {
	// Quota is the current quota status after the action.
	Quota model.Quota `json:"quota"`
	// Stats is the latest statistics for the image.
	Stats statModel `json:"stats"`
}
