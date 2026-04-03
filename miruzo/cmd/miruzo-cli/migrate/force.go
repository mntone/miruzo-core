package migrate

import (
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/service/migration"
	"github.com/spf13/cobra"
)

var forceCommand = &cobra.Command{
	Use:   "force V",
	Short: "Set version V but don't run migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		return withMigrationService(command, func(srv migration.Service) error {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return srv.SetVersion(command.Context(), version)
		})
	},
}

func init() {
	Command.AddCommand(forceCommand)
}
