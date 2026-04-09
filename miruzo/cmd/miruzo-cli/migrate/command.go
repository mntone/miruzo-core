package migrate

import (
	"errors"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	"github.com/mntone/miruzo-core/miruzo/internal/service/migration"
	"github.com/spf13/cobra"
)

type migrationCommandCallback func(srv migration.Service) error

func withMigrationService(
	command *cobra.Command,
	callback migrationCommandCallback,
) (err error) {
	cfg, err := app.LoadConfig()
	if err != nil {
		return err
	}

	hdl, err := persistence.OpenManagementHandle(command.Context(), cfg.Database)
	if err != nil {
		return err
	}

	srv, err := migration.New(hdl.MigrationRunner())
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, hdl.Close())
	}()

	return callback(srv)
}

var Command = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}
