package reaction

import (
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
)

func mapLoveResponse(r reaction.LoveResult) loveResponse {
	return loveResponse{
		Quota: r.Quota,
		Stats: statModel{
			Score:        r.Stats.Score,
			FirstLovedAt: r.Stats.FirstLovedAt.ToPointer(),
			LastLovedAt:  r.Stats.LastLovedAt.ToPointer(),
		},
	}
}
