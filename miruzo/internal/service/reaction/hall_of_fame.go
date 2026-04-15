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
	hallOfFameGrantedAt := srv.clk.Now()
	periodStartAt := srv.dailyPeriodResolver.PeriodStart(hallOfFameGrantedAt)

	err := srv.prov.Session(requestContext, func(ctx context.Context, repos persist.SessionRepositories) error {
		err := repos.Stats().ApplyHallOfFameGranted(
			ctx,
			ingestID,
			hallOfFameGrantedAt,
			srv.hallOfFameScoreThreshold,
		)
		if err != nil {
			return err
		}

		err = repos.Action().CreateHallOfFameIfAbsent(
			ctx,
			ingestID,
			persist.HallOfFameActionTypeGranted,
			hallOfFameGrantedAt,
			periodStartAt,
		)
		return err
	})
	if err != nil {
		return HallOfFameResult{}, serviceerror.MapPersistError(err)
	}

	return HallOfFameResult{
		Stats: model.HallOfFameStats{
			HallOfFameAt: mo.Some(hallOfFameGrantedAt),
		},
	}, nil
}

func (srv *Service) RevokeHallOfFame(
	requestContext context.Context,
	ingestID model.IngestIDType,
) (HallOfFameResult, error) {
	hallOfFameRevokedAt := srv.clk.Now()
	periodStartAt := srv.dailyPeriodResolver.PeriodStart(hallOfFameRevokedAt)

	err := srv.prov.Session(requestContext, func(ctx context.Context, repos persist.SessionRepositories) error {
		err := repos.Stats().ApplyHallOfFameRevoked(ctx, ingestID)
		if err != nil {
			return err
		}

		err = repos.Action().CreateHallOfFameIfAbsent(
			ctx,
			ingestID,
			persist.HallOfFameActionTypeRevoked,
			hallOfFameRevokedAt,
			periodStartAt,
		)
		return err
	})
	if err != nil {
		return HallOfFameResult{}, serviceerror.MapPersistError(err)
	}

	return HallOfFameResult{}, nil
}
