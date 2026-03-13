package quota

import (
	apiModel "github.com/mntone/miruzo-core/miruzo/internal/api/model"
	"github.com/mntone/miruzo-core/miruzo/internal/service/user"
)

func mapQuota(result user.QuotaResult) (quotaResponse, error) {
	loveQuota, err := apiModel.MapQuota(result.Love)
	if err != nil {
		return quotaResponse{}, err
	}

	return quotaResponse{
		Love: loveQuota,
	}, nil
}
