package quota

import (
	"github.com/mntone/miruzo-core/miruzo/internal/service/user"
)

func mapQuota(result user.QuotaResult) quotaResponse {
	return quotaResponse{
		Love: result.Love,
	}
}
