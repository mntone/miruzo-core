package reaction

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

func (srv *Service) Love(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (LoveResult, error) {
	lovedAt := srv.clk.Now()
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
	err := srv.prov.Session(requestContext, func(ctx context.Context, repos persist.SessionRepositories) error {
		dailyLoveUsed, err := repos.User().IncrementDailyLoveUsed(ctx, srv.dailyLoveLimit)
		if err != nil {
			return err
		}
		if dailyLoveUsed < srv.dailyLoveLimit {
			result.Quota.Remaining = srv.dailyLoveLimit - dailyLoveUsed
		}

		stats, err := repos.Stats().ApplyLove(
			ctx,
			ingestID,
			scoreDelta,
			lovedAt,
			srv.hallOfFameScoreThreshold,
			periodStartAt,
		)
		if err != nil {
			return err
		}

		err = repos.Action().CreateLoveIfAbsent(
			ctx,
			ingestID,
			persist.LoveActionTypeLove,
			lovedAt,
			periodStartAt,
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
