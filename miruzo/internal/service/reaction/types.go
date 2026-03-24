package reaction

import "github.com/mntone/miruzo-core/miruzo/internal/model"

type LoveResult struct {
	Quota model.Quota
	Stats model.LoveStats
}

type HallOfFameResult struct {
	Stats model.HallOfFameStats
}
