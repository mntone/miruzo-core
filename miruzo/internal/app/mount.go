package app

import (
	"net/http"

	imageListAPI "github.com/mntone/miruzo-core/miruzo/internal/api/image/list"
	"github.com/mntone/miruzo-core/miruzo/internal/api/variant"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	imageListService "github.com/mntone/miruzo-core/miruzo/internal/service/imagelist"
)

func MountAPI(
	mux *http.ServeMux,
	factory persist.RepositoryFactory,
	cfg config.AppConfig,
) {
	imageListBackoff := newBackoffPolicyFromConfig(cfg.API.Retry.Read)
	imageListService := imageListService.New(
		factory.NewImageList(),
		imageListBackoff,
		cfg.Score.EngagedScoreThreshold,
	)
	imageListHandler := imageListAPI.NewHandler(
		imageListService,
		cfg.API.VariantLayers,
		variant.NewMediaURLBuilder(cfg.API.MediaPublic),
	)
	imageListAPI.RegisterRoutes(mux, imageListHandler)
}
