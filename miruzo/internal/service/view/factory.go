package view

import (
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/retry/backoff"
)

type Service struct {
	mgr             persist.PersistenceManager
	backoff         backoff.Policy
	scoreCalculator score.Calculator
	viewMilestones  []int64
}

func New(
	persistenceManager persist.PersistenceManager,
	backoff backoff.Policy,
	scoreCalculator score.Calculator,
	viewMilestones []int64,
) Service {
	return Service{
		mgr:             persistenceManager,
		backoff:         backoff,
		scoreCalculator: scoreCalculator,
		viewMilestones:  viewMilestones,
	}
}
