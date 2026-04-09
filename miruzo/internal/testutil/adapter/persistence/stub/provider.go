package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type PersistenceProvider struct {
	session
}

func NewStubPersistenceProvider(
	dailyLoveUsed int32,
	statsEntries ...model.Stats,
) *PersistenceProvider {
	return &PersistenceProvider{
		session: session{
			ActionStub: NewStubActionRepository(),
			JobStub:    NewStubJobRepository(),
			StatsStub:  NewStubStatsRepository(statsEntries...),
			UserStub:   NewStubUserRepository(dailyLoveUsed),
		},
	}
}

func (prov *PersistenceProvider) Close() error {
	return nil
}

func (prov *PersistenceProvider) Repos() persist.Repositories {
	return &prov.session
}

func (prov *PersistenceProvider) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	// Get snapshots
	a := prov.ActionStub.snapshot()
	j := prov.JobStub.snapshot()
	s := prov.StatsStub.snapshot()
	u := prov.UserStub.snapshot()

	err := callback(ctx, &txSession{session: &prov.session})
	if err != nil {
		// Rollback
		prov.ActionStub.actionStorage = a
		prov.JobStub.jobStorage = j
		prov.StatsStub.statsStorage = s
		prov.UserStub.userStorage = u
	}

	return err
}
