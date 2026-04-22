package database

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create application database",
	Args:  cobra.NoArgs,
	RunE: func(command *cobra.Command, args []string) error {
		return withDatabaseAdminHandle(
			command,
			func(hdl persistence.DatabaseAdminHandle) error {
				exists, err := hdl.Exists(command.Context())
				if err != nil {
					return fmt.Errorf("check database exists: %w", err)
				}
				if exists {
					return fmt.Errorf("database already exists")
				}

				if err := hdl.Create(command.Context()); err != nil {
					return err
				}
				return nil
			},
		)
	},
}

func init() {
	Command.AddCommand(createCommand)
}
