package reaction

import (
	apiModel "github.com/mntone/miruzo-core/miruzo/internal/api/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
)

func mapLoveResponse(r reaction.LoveResult) (loveResponse, error) {
	quota, err := apiModel.MapQuota(r.Quota)
	if err != nil {
		return loveResponse{}, err
	}

	return loveResponse{
		Quota: quota,
		Stats: statModel{
			Score:        r.Stats.Score,
			FirstLovedAt: r.Stats.FirstLovedAt.ToPointer(),
			LastLovedAt:  r.Stats.LastLovedAt.ToPointer(),
		},
	}, nil
}
