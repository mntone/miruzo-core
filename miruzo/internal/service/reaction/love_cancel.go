package reaction

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

func (srv *Service) LoveCancel(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (LoveResult, error) {
	canceledAt := srv.clk.Now()
	periodStartAt, periodEndAt := srv.dailyPeriodResolver.PeriodRange(canceledAt)
	scoreDelta := srv.scoreCalculator.LoveCanceledDelta()

	result := LoveResult{
		Quota: model.Quota{
			Period:    model.PeriodTypeDaily,
			ResetAt:   periodEndAt,
			Limit:     srv.dailyLoveLimit,
			Remaining: 0,
		},
	}
	err := srv.prov.Session(requestContext, func(ctx context.Context, repos persist.SessionRepositories) error {
		stats, err := repos.Stats().ApplyLoveCanceled(ctx, ingestID, scoreDelta, canceledAt, periodStartAt)
		if err != nil {
			return err
		}

		err = repos.Action().CreateLoveIfAbsent(
			ctx,
			ingestID,
			persist.LoveActionTypeLoveCanceled,
			canceledAt,
			periodStartAt,
		)
		if err != nil {
			return err
		}

		dailyLoveUsed, err := repos.User().DecrementDailyLoveUsed(ctx)
		if err != nil {
			if !errors.Is(err, persist.ErrQuotaUnderflow) {
				return err
			}

			log.Printf(
				"daily love used underflow: ingest_id=%d love_canceled_at=%s",
				ingestID,
				canceledAt.Format(time.RFC3339Nano),
			)
		}

		result.Quota.Remaining = srv.dailyLoveLimit - dailyLoveUsed
		result.Stats = stats
		return nil
	})
	if err != nil {
		return result, serviceerror.MapPersistError(err)
	}

	return result, nil
}
