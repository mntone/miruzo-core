package migrate

import (
	"fmt"
	"strconv"

	"github.com/mntone/miruzo-core/miruzo/internal/service/migration"
	"github.com/spf13/cobra"
)

var upCommand = &cobra.Command{
	Use:   "up [N]",
	Short: "Apply pending migrations",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(command *cobra.Command, args []string) error {
		return withMigrationService(command, func(srv migration.Service) error {
			if len(args) == 0 {
				return srv.Up(command.Context())
			}

			steps, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if steps <= 0 {
				return fmt.Errorf("N must be greater than 0")
			}

			return srv.Step(command.Context(), steps)
		})
	},
}

func init() {
	Command.AddCommand(upCommand)
}
