package imagelist

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

type Service struct {
	repository            persist.ImageListRepository
	backoff               backoff.Policy
	engagedScoreThreshold model.ScoreType
	variantLayersBuilder  *media.VariantLayersBuilder
}

func New(
	repo persist.ImageListRepository,
	backoff backoff.Policy,
	engagedScoreThreshold model.ScoreType,
	variantLayersBuilder *media.VariantLayersBuilder,
) (Service, error) {
	if variantLayersBuilder == nil {
		return Service{}, fmt.Errorf("variantLayersBuilder must not be nil")
	}

	return Service{
		repository:            repo,
		backoff:               backoff,
		engagedScoreThreshold: engagedScoreThreshold,
		variantLayersBuilder:  variantLayersBuilder,
	}, nil
}
