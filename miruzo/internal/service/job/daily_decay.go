package job

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
)

const applyDailyDecayBatchCount int32 = 500

type DailyDecayService struct {
	prov                persist.PersistenceProvider
	clk                 clock.Provider
	dailyPeriodResolver period.DailyResolver
	scoreCalculator     score.Calculator
}

func NewDailyDecay(
	persistenceProvider persist.PersistenceProvider,
	clockProvider clock.Provider,
	dailyPeriodResolver period.DailyResolver,
	scoreCalculator score.Calculator,
) *DailyDecayService {
	return &DailyDecayService{
		prov:                persistenceProvider,
		clk:                 clockProvider,
		dailyPeriodResolver: dailyPeriodResolver,
		scoreCalculator:     scoreCalculator,
	}
}

func (srv *DailyDecayService) ApplyDailyDecay(ctx context.Context) error {
	occurredAt := srv.clk.Now()
	periodStartAt := srv.dailyPeriodResolver.PeriodStart(occurredAt)

	err := srv.prov.Session(ctx, func(ctx context.Context, repos persist.SessionRepositories) error {
		reporter := dailyDecayProgressReporter{}
		defer reporter.Print()

		iter := repos.Stats().IterateStatsForDailyDecay(ctx, applyDailyDecayBatchCount)
		for row, err := range iter {
			if err != nil {
				reporter.AddFailed(err)
				return err
			}

			lastViewedAt, present := row.LastViewedAt.Get()
			if !present {
				reporter.AddSkipped()
				continue
			}

			err = repos.Action().CreateDailyDecayIfAbsent(
				ctx,
				row.IngestID,
				occurredAt,
				periodStartAt,
			)
			if err != nil {
				if errors.Is(err, persist.ErrConflict) {
					reporter.AddSkipped()
					continue
				}

				reporter.AddFailed(err)
				return err
			}

			newScore, negative := srv.scoreCalculator.DailyDecay(
				row.Score,
				lastViewedAt,
				occurredAt,
			)
			if negative {
				log.Printf(
					"daily decay negative: ingest_id=%d last_viewed_at=%s occurred_at=%s",
					row.IngestID,
					lastViewedAt.Format(time.RFC3339Nano),
					occurredAt.Format(time.RFC3339Nano),
				)
			}

			err = repos.Stats().ApplyDecay(
				ctx,
				row.IngestID,
				newScore,
				occurredAt,
			)
			if err != nil {
				reporter.AddFailed(err)
				return err
			}

			reporter.AddProcessed()
		}

		return repos.User().ResetDailyLoveUsed(ctx)
	})
	if err != nil {
		return serviceerror.MapPersistError(err)
	}

	return nil
}
