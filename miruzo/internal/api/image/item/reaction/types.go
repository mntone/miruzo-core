package reaction

import (
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
)

type loveStatModel struct {
	Score        model.ScoreType `json:"score"`
	FirstLovedAt *time.Time      `json:"first_loved_at"`
	LastLovedAt  *time.Time      `json:"last_loved_at"`
}

// Response payload for love actions.
type loveResponse struct {
	// Quota is the current quota status after the action.
	Quota model.Quota `json:"quota"`
	// Stats is the latest statistics for the image.
	Stats loveStatModel `json:"stats"`
}

type hallOfFameStatModel struct {
	HallOfFameAt *time.Time `json:"hall_of_fame_at"`
}

// Response payload for hall of fame actions.
type hallOfFameResponse struct {
	// Stats is the latest statistics for the image.
	Stats hallOfFameStatModel `json:"stats"`
}
