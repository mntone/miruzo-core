package view

import (
	"context"
	"log"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/samber/mo"
)

type ContextArgs struct {
	IngestID model.IngestIDType

	// ExcludeFormats lists variant formats to exclude.
	// Nil keeps the default format policy.
	ExcludeFormats []media.ImageFormat
}

func (srv *Service) shouldTriggerViewMilestone(stats model.Stats) bool {
	for _, milestone := range srv.viewMilestones {
		if stats.ViewCount >= milestone && stats.ViewMilestoneCount < milestone {
			return true
		}
	}
	return false
}

func (srv *Service) GetContext(requestContext context.Context, args ContextArgs) (model.ImageWithStats, error) {
	viewedAt := srv.clk.Now()

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
					return repos.View.GetImageWithStatsForUpdate(requestContext, args.IngestID)
				},
			)
			if err != nil {
				return err
			}

			scoreDelta, negative := srv.scoreCalculator.ViewDelta(imageWithStats.Stats.LastViewedAt, viewedAt)
			if negative {
				lastViewedAt := imageWithStats.Stats.LastViewedAt.MustGet()
				log.Printf(
					"view delta negative: ingest_id=%d last_viewed_at=%s viewed_at=%s",
					args.IngestID,
					lastViewedAt.Format(time.RFC3339Nano),
					viewedAt.Format(time.RFC3339Nano),
				)
			}

			if srv.shouldTriggerViewMilestone(imageWithStats.Stats) {
				err = repos.Stats.ApplyViewWithMilestone(
					ctx,
					args.IngestID,
					scoreDelta,
					viewedAt,
				)

				imageWithStats.Stats.LastViewedAt = mo.Some(viewedAt)
				imageWithStats.Stats.ViewCount += 1
				imageWithStats.Stats.ViewMilestoneCount = imageWithStats.Stats.ViewCount
			} else {
				err = repos.Stats.ApplyView(
					ctx,
					args.IngestID,
					scoreDelta,
					viewedAt,
				)

				imageWithStats.Stats.LastViewedAt = mo.Some(viewedAt)
				imageWithStats.Stats.ViewCount += 1
			}
			if err != nil {
				return err
			}

			_, err = repos.Action.Create(
				ctx,
				args.IngestID,
				model.ActionTypeView,
				viewedAt,
			)
			if err != nil {
				return err
			}

			result = imageWithStats
			return nil
		},
	)
	if err != nil {
		return model.ImageWithStats{}, serviceerror.MapPersistError(err)
	}

	options := media.VariantFilterOptions{
		IncludeFormatSet: media.ComputeAllowedFormats(args.ExcludeFormats),
		KeepFallback:     true,
	}
	variants := result.Layers.ToDomain().FilterWith(options)
	layers := srv.variantLayersBuilder.GroupVariantsByLayer(variants)
	return result.ToDTO(layers), nil
}
