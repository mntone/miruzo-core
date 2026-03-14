package reaction

import (
	"context"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

type LoveResult struct {
	Quota model.Quota
	Stats persist.LoveStats
}

func (srv Service) Love(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (LoveResult, error) {
	lovedAt := time.Now().UTC()
	periodStartAt, periodEndAt := srv.dailyPeriodResolver.PeriodRange(lovedAt)
	scoreDelta := srv.scoreCalculator.LoveDelta()

	result := LoveResult{
		Quota: model.Quota{
			Period:    model.PeriodTypeDaily,
			ResetAt:   periodEndAt,
			Limit:     srv.dailyLoveLimit,
			Remaining: 0,
		},
	}
	err := srv.mgr.Session(requestContext, func(ctx context.Context, repos persist.Repositories) error {
		dailyLoveUsed, err := repos.User.IncrementDailyLoveUsed(ctx, srv.dailyLoveLimit)
		if err != nil {
			return err
		}
		if dailyLoveUsed < srv.dailyLoveLimit {
			result.Quota.Remaining = srv.dailyLoveLimit - dailyLoveUsed
		}

		stats, err := repos.Stats.ApplyLove(ctx, ingestID, scoreDelta, lovedAt, periodStartAt)
		if err != nil {
			return err
		}

		_, err = repos.Action.Create(
			ctx,
			ingestID,
			model.ActionTypeLove,
			lovedAt,
		)
		if err != nil {
			return err
		}

		result.Stats = stats
		return nil
	})
	if err != nil {
		return result, serviceerror.MapPersistError(err)
	}

	return result, nil
}
