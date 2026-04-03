package migrate

import (
	"fmt"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/service/migration"
	"github.com/spf13/cobra"
)

var gotoCommand = &cobra.Command{
	Use:   "goto V",
	Short: "Migrate to version V",
	Args:  cobra.ExactArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		return withMigrationService(command, func(srv migration.Service) error {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if version < 0 {
				return fmt.Errorf("V must be greater than or equal to 0")
			}

			return srv.Migrate(command.Context(), version)
		})
	},
}

func init() {
	Command.AddCommand(gotoCommand)
}
