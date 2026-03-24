package reaction

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/serviceerror"
	"github.com/samber/mo"
)

func (srv *Service) GrantHallOfFame(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (HallOfFameResult, error) {
	hallOfFameAt := srv.clk.Now()

	err := srv.mgr.Session(requestContext, func(ctx context.Context, repos persist.Repositories) error {
		err := repos.Stats.ApplyHallOfFameGranted(
			ctx,
			ingestID,
			hallOfFameAt,
			srv.hallOfFameScoreThreshold,
		)
		if err != nil {
			return err
		}

		_, err = repos.Action.Create(
			ctx,
			ingestID,
			model.ActionTypeHallOfFameGranted,
			hallOfFameAt,
		)
		return err
	})
	if err != nil {
		return HallOfFameResult{}, serviceerror.MapPersistError(err)
	}

	return HallOfFameResult{
		Stats: model.HallOfFameStats{
			HallOfFameAt: mo.Some(hallOfFameAt),
		},
	}, nil
}

func (srv *Service) RevokeHallOfFame(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (HallOfFameResult, error) {
	hallOfFameAt := srv.clk.Now()

	err := srv.mgr.Session(requestContext, func(ctx context.Context, repos persist.Repositories) error {
		err := repos.Stats.ApplyHallOfFameRevoked(ctx, ingestID)
		if err != nil {
			return err
		}

		_, err = repos.Action.Create(
			ctx,
			ingestID,
			model.ActionTypeHallOfFameRevoked,
			hallOfFameAt,
		)
		return err
	})
	if err != nil {
		return HallOfFameResult{}, serviceerror.MapPersistError(err)
	}

	return HallOfFameResult{}, nil
}
