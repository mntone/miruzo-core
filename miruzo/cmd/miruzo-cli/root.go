package main

import (
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/cmd/miruzo-cli/database"
	"github.com/mntone/miruzo-core/miruzo/cmd/miruzo-cli/job"
	"github.com/mntone/miruzo-core/miruzo/cmd/miruzo-cli/migrate"
	"github.com/spf13/cobra"
)

var showVersion bool

var rootCommand = &cobra.Command{
	Use:   "miruzo-cli",
	Short: "Miruzo command line tools",
	Args:  cobra.NoArgs,
	RunE: func(command *cobra.Command, args []string) error {
		if showVersion {
			_, err := fmt.Fprintln(command.OutOrStdout(), version)
			return err
		}

		return command.Help()
	},
}

func init() {
	rootCommand.Flags().BoolVarP(
		&showVersion,
		"version",
		"v",
		false,
		"Print CLI version",
	)

	rootCommand.AddCommand(database.Command)
	rootCommand.AddCommand(job.Command)
	rootCommand.AddCommand(migrate.Command)
}
