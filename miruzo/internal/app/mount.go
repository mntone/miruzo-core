package app

import (
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api"
	healthAPI "github.com/mntone/miruzo-core/miruzo/internal/api/health"
	contextAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/item/context"
	reactionAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/item/reaction"
	imageListAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/list"
	quotaAPI "github.com/mntone/miruzo-core/miruzo/internal/api/quota"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
	reactionService "github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
	userService "github.com/mntone/miruzo-core/miruzo/internal/service/user"
	viewService "github.com/mntone/miruzo-core/miruzo/internal/service/view"
)

func buildScoreCalculator(
	dailyResolver period.DailyResolver,
	cfg config.ScoreConfig,
) score.Calculator {
	return score.New(
		dailyResolver,
		cfg.ViewBonusAtFirst,
		cfg.ViewBonusByDays,
		cfg.ViewBonusFallback,
		cfg.MemoBonus,
		cfg.MemoPenalty,
		cfg.LoveBonus,
		cfg.LovePenalty,
	)
}

func MountAPI(
	mux *http.ServeMux,
	manager persist.PersistenceManager,
	cfg config.AppConfig,
	version string,
) {
	readBackoff := newBackoffPolicyFromConfig(cfg.API.Retry.Read)
	imageListService := imageListService.New(
		manager.Repos().ImageList,
		readBackoff,
		cfg.Score.EngagedScoreThreshold,
	)
	mediaURLBuilder := variant.NewMediaURLBuilder(cfg.API.MediaPublic)
	imageListHandler := imageListAPI.NewHandler(
		imageListService,
		cfg.API.VariantLayers,
		mediaURLBuilder,
	)
	imageListAPI.RegisterRoutes(mux, imageListHandler)

	dailyResolver := period.NewDailyResolver(cfg.Period.DayStartOffset)
	scoreCalculator := buildScoreCalculator(dailyResolver, cfg.Score)
	viewService := viewService.New(manager, readBackoff, scoreCalculator, cfg.View.Milestones)
	imageItemHandler := contextAPI.NewHandler(viewService, cfg.API.VariantLayers, mediaURLBuilder)
	contextAPI.RegisterRoutes(mux, imageItemHandler)

	reactionService := reactionService.New(
		manager,
		dailyResolver,
		scoreCalculator,
		cfg.Quota.DailyLoveLimit,
	)
	reactionHandler := reactionAPI.NewHandler(reactionService)
	reactionAPI.RegisterRoutes(mux, reactionHandler)

	userService := userService.New(
		manager.Repos().User,
		dailyResolver,
		cfg.Quota.DailyLoveLimit,
	)
	quotaHandler := quotaAPI.NewHandler(userService)
	quotaAPI.RegisterRoutes(mux, quotaHandler)

	healthHandler := healthAPI.NewHandler(version)
	healthAPI.RegisterRoutes(mux, healthHandler)

	api.RegisterNotFoundRoute(mux)
}
