package app

import (
	"context"
	"log"
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/api"
	healthAPI "github.com/mntone/miruzo-core/miruzo/internal/api/health"
	contextAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/item/context"
	reactionAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/item/reaction"
	imageListAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/list"
	"github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
	quotaAPI "github.com/mntone/miruzo-core/miruzo/internal/api/quota"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
	reactionService "github.com/mntone/miruzo-core/miruzo/internal/service/reaction"
	userService "github.com/mntone/miruzo-core/miruzo/internal/service/user"
	viewService "github.com/mntone/miruzo-core/miruzo/internal/service/view"
)

func BuildScoreCalculator(
	dailyResolver period.DailyResolver,
	cfg config.ScoreConfig,
) score.Calculator {
	return score.New(
		dailyResolver,
		cfg.DailyDecayNoAccessAdjustment,
		cfg.DailyDecayPenalty,
		cfg.DailyDecayInterval10dPenalty,
		cfg.DailyDecayHighScorePenalty,
		cfg.DailyDecayHighScoreThreshold,
		cfg.ViewBonusAtFirst,
		cfg.ViewBonusByDays,
		cfg.ViewBonusFallback,
		cfg.MemoBonus,
		cfg.MemoPenalty,
		cfg.LoveBonus,
		cfg.LovePenalty,
	)
}

func mountAPI(
	mux *http.ServeMux,
	manager persist.PersistenceManager,
	cfg config.AppConfig,
	version string,
) {
	dailyResolver, err := NewDailyResolver(
		context.Background(),
		cfg.Period,
		manager.Repos().Settings,
	)
	if err != nil {
		log.Fatalf("app: failed to build daily resolver: %v", err)
	}

	cors := middleware.NewCORSFactory(cfg.CORS.AllowOrigins, cfg.CORS.MaxAge)
	variantLayersBuilder := media.NewVariantLayerBuilder(cfg.API.VariantLayers)
	readBackoff := newBackoffPolicyFromConfig(cfg.API.Retry.Read)
	imageListService, err := imageListService.New(
		manager.Repos().ImageList,
		readBackoff,
		cfg.Score.EngagedScoreThreshold,
		variantLayersBuilder,
	)
	if err != nil {
		log.Fatalf("app: failed to build image list service: %v", err)
	}
	mediaURLBuilder := variant.NewMediaURLBuilder(cfg.API.MediaPublic)
	imageListHandler := imageListAPI.NewHandler(
		imageListService,
		cfg.API.VariantLayers,
		mediaURLBuilder,
	)
	imageListAPI.RegisterRoutes(mux, cors, imageListHandler)

	clockProvider := clock.NewSystemProvider()
	scoreCalculator := BuildScoreCalculator(dailyResolver, cfg.Score)
	viewService, err := viewService.New(
		manager,
		readBackoff,
		clockProvider,
		scoreCalculator,
		variantLayersBuilder,
		cfg.View.Milestones,
	)
	if err != nil {
		log.Fatalf("app: failed to build view service: %v", err)
	}
	imageItemHandler := contextAPI.NewHandler(viewService, mediaURLBuilder)
	contextAPI.RegisterRoutes(mux, cors, imageItemHandler)

	reactionService, err := reactionService.New(
		manager,
		clockProvider,
		dailyResolver,
		scoreCalculator,
		cfg.Quota.DailyLoveLimit,
		cfg.Score.HallOfFameScoreThreshold,
	)
	if err != nil {
		log.Fatalf("app: failed to build reaction service: %v", err)
	}
	reactionHandler := reactionAPI.NewHandler(reactionService)
	reactionAPI.RegisterRoutes(mux, cors, reactionHandler)

	userService, err := userService.New(
		manager.Repos().User,
		clockProvider,
		dailyResolver,
		cfg.Quota.DailyLoveLimit,
	)
	if err != nil {
		log.Fatalf("app: failed to build user service: %v", err)
	}
	quotaHandler := quotaAPI.NewHandler(userService)
	quotaAPI.RegisterRoutes(mux, cors, quotaHandler)

	healthHandler := healthAPI.NewHandler(version)
	healthAPI.RegisterRoutes(mux, cors, healthHandler)

	api.RegisterNotFoundRoute(mux)

}

func MountAll(
	mux *http.ServeMux,
	manager persist.PersistenceManager,
	cfg config.AppConfig,
	version string,
) {
	mountAPI(mux, manager, cfg, version)

	if cfg.Server.StaticFiles.Enabled {
		mountStatic(
			mux,
			cfg.Server.StaticFiles.RootDirectory,
			cfg.API.MediaPublic.BasePath,
			buildCacheControlHeader(cfg.Server.StaticFiles),
			cfg.Server.StaticFiles.NoSniff,
		)
	}
}
