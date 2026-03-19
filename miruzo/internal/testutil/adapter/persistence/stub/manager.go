package stub

import (
	"context"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
)

type PersistenceManager struct {
	Action       *actionRepository
	Stats        *statsRepository
	User         *UserRepository
	Repositories persist.Repositories
}

func NewStubPersistenceManager(
	dailyLoveUsed int32,
	statsEntries ...persist.Stats,
) PersistenceManager {
	action := NewStubActionRepository()
	stats := NewStubStatsRepository(statsEntries...)
	user := NewStubUserRepository(dailyLoveUsed)
	return PersistenceManager{
		Action: action,
		Stats:  stats,
		User:   user,
		Repositories: persist.Repositories{
			Action:    action,
			ImageList: nil,
			Settings:  nil,
			Stats:     stats,
			User:      user,
			View:      nil,
		},
	}
}

func (mgr PersistenceManager) Close() error {
	return nil
}

func (mgr PersistenceManager) Repos() persist.Repositories {
	return mgr.Repositories
}

func (mgr PersistenceManager) Session(
	ctx context.Context,
	callback persist.SessionCallback,
) error {
	// Get snapshots
	a := mgr.Action.snapshot()
	s := mgr.Stats.snapshot()
	u := mgr.User.snapshot()

	err := callback(ctx, mgr.Repositories)
	if err != nil {
		// Rollback
		mgr.Action.actionStorage = a
		mgr.Stats.statsStorage = s
		mgr.User.userStorage = u
	}

	return err
}
