package view

import (
	"context"
	"log"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry"
	"github.com/samber/mo"
)

func (srv Service) shouldTriggerViewMilestone(stats persist.Stats) bool {
	for _, milestone := range srv.viewMilestones {
		if stats.ViewCount >= milestone && stats.ViewMilestoneCount < milestone {
			return true
		}
	}
	return false
}

func (srv Service) GetContext(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (persist.ImageWithStats, error) {
	evaluatedAt := time.Now().UTC()

	var result persist.ImageWithStats
	err := srv.mgr.Session(
		requestContext,
		func(
			ctx context.Context,
			repos persist.Repositories,
		) error {
			imageWithStats, err := retry.Retry(
				ctx,
				srv.backoff,
				func(requestContext context.Context) (persist.ImageWithStats, error) {
					return repos.View.GetImageWithStats(requestContext, ingestID)
				},
			)
			if err != nil {
				return err
			}

			scoreDelta, negative := srv.scoreCalculator.ViewDelta(imageWithStats.Stats.LastViewedAt, evaluatedAt)
			if negative {
				lastViewedAt := imageWithStats.Stats.LastViewedAt.MustGet()
				log.Printf(
					"view delta negative: ingest_id=%d last_viewed_at=%s evaluated_at=%s",
					ingestID,
					lastViewedAt.Format(time.RFC3339Nano),
					evaluatedAt.Format(time.RFC3339Nano),
				)
			}

			if srv.shouldTriggerViewMilestone(imageWithStats.Stats) {
				err = repos.Stats.ApplyViewWithMilestone(
					ctx,
					ingestID,
					scoreDelta,
					evaluatedAt,
				)

				imageWithStats.Stats.LastViewedAt = mo.Some(evaluatedAt)
				imageWithStats.Stats.ViewCount += 1
				imageWithStats.Stats.ViewMilestoneCount = imageWithStats.Stats.ViewCount
			} else {
				err = repos.Stats.ApplyView(
					ctx,
					ingestID,
					scoreDelta,
					evaluatedAt,
				)

				imageWithStats.Stats.LastViewedAt = mo.Some(evaluatedAt)
				imageWithStats.Stats.ViewCount += 1
			}
			if err != nil {
				return err
			}

			_, err = repos.Action.CreateAction(
				ctx,
				ingestID,
				model.ActionTypeView,
				evaluatedAt,
			)
			if err != nil {
				return err
			}

			result = imageWithStats
			return nil
		},
	)
	if err != nil {
		return persist.ImageWithStats{}, err
	}

	return result, nil
}
