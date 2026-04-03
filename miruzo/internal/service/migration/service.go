package migration

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type Service struct {
	persist.MigrationRunner
}

func New(runner persist.MigrationRunner) (Service, error) {
	if runner == nil {
		return Service{}, fmt.Errorf("runner must not be nil")
	}

	return Service{
		MigrationRunner: runner,
	}, nil
}
