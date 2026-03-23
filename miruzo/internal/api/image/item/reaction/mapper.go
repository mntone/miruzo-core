package reaction

import (
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
)

func mapLoveResponse(r reaction.LoveResult) loveResponse {
	return loveResponse{
		Quota: r.Quota,
		Stats: loveStatModel{
			Score:        r.Stats.Score,
			FirstLovedAt: r.Stats.FirstLovedAt.ToPointer(),
			LastLovedAt:  r.Stats.LastLovedAt.ToPointer(),
		},
	}
}

func mapHallOfFameResponse(r reaction.HallOfFameResult) hallOfFameResponse {
	return hallOfFameResponse{
		Stats: hallOfFameStatModel{
			HallOfFameAt: r.Stats.HallOfFameAt.ToPointer(),
		},
	}
}
