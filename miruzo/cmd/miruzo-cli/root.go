package main

import (
	"github.com/mntone/miruzo-core/miruzo/cmd/miruzo-cli/job"
	"github.com/mntone/miruzo-core/miruzo/cmd/miruzo-cli/migrate"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "miruzo-cli",
	Short: "Miruzo command line tools",
}

func init() {
	rootCommand.AddCommand(job.Command)
	rootCommand.AddCommand(migrate.Command)
}
