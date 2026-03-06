package imagelist

import (
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

type Service struct {
	repository            persist.ImageListRepository
	backoff               backoff.Policy
	engagedScoreThreshold int16
}

func New(
	repo persist.ImageListRepository,
	backoff backoff.Policy,
	engagedScoreThreshold int16,
) Service {
	return Service{
		repository:            repo,
		backoff:               backoff,
		engagedScoreThreshold: engagedScoreThreshold,
	}
}
