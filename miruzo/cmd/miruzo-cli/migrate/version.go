package migrate

import (
	"github.com/mntone/miruzo-core/miruzo/internal/service/migration"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print current migration version",
	RunE: func(command *cobra.Command, args []string) error {
		return withMigrationService(command, func(srv migration.Service) error {
			version, err := srv.Version(command.Context())
			if err != nil {
				return err
			}

			command.Println(version)
			return nil
		})
	},
}

func init() {
	Command.AddCommand(versionCommand)
}
