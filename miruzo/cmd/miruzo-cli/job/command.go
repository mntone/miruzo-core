package job

import (
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/jobguard"
	"github.com/spf13/cobra"
)

type jobCommandCallback func(cfg config.AppConfig, mgr persist.PersistenceProvider) error

func withJobRunGuard(
	name string,
	command *cobra.Command,
	callback jobCommandCallback,
) (err error) {
	cfg, err := app.LoadConfig()
	if err != nil {
		return err
	}

	hdl, err := persistence.OpenManagementHandle(command.Context(), cfg.Database)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, hdl.Close())
	}()

	clk := clock.NewSystemProvider()
	prov := hdl.PersistenceProvider()
	guard := jobguard.NewWithJobRepository(prov.Repos().Job())
	acquired, err := guard.TryAcquire(command.Context(), name, clk.Now())
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}
	defer func() {
		releaseErr := guard.Release(command.Context(), name, clk.Now())
		if releaseErr != nil {
			err = errors.Join(err, releaseErr)
		}
	}()

	return callback(cfg, prov)
}

var Command = &cobra.Command{
	Use:   "job",
	Short: "Run maintenance jobs",
}
