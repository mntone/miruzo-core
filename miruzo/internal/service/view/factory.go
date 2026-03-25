package view

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/media"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

type Service struct {
	mgr                  persist.PersistenceManager
	backoff              backoff.Policy
	clk                  clock.Provider
	scoreCalculator      score.Calculator
	variantLayersBuilder *media.VariantLayersBuilder
	viewMilestones       []int64
}

func New(
	persistenceManager persist.PersistenceManager,
	backoff backoff.Policy,
	clockProvider clock.Provider,
	scoreCalculator score.Calculator,
	variantLayersBuilder *media.VariantLayersBuilder,
	viewMilestones []int64,
) (Service, error) {
	if variantLayersBuilder == nil {
		return Service{}, fmt.Errorf("variantLayersBuilder must not be nil")
	}

	return Service{
		mgr:                  persistenceManager,
		backoff:              backoff,
		clk:                  clockProvider,
		scoreCalculator:      scoreCalculator,
		variantLayersBuilder: variantLayersBuilder,
		viewMilestones:       viewMilestones,
	}, nil
}
