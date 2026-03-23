package reaction

import (
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type LoveResult struct {
	Quota model.Quota
	Stats persist.LoveStats
}

type HallOfFameResult struct {
	Stats persist.HallOfFameStats
}
