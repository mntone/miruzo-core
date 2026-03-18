package persist

import "github.com/mntone/miruzo-core/miruzo/internal/model"

type User struct {
	ID            int16
	DailyLoveUsed model.QuotaInt
}
