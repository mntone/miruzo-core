package job

import (
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/clock"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	jobService "github.com/mntone/miruzo-core/miruzo/internal/service/job"
	"github.com/spf13/cobra"
)

var dailyDecayCommand = &cobra.Command{
	Use:   "daily-decay",
	Short: "Apply daily score decay",
	RunE: func(command *cobra.Command, args []string) error {
		return withJobRunGuard("daily_decay", command, func(cfg config.AppConfig, mgr persist.PersistenceManager) error {
			dailyPeriodResolver, err := app.NewDailyResolver(
				command.Context(),
				cfg.Period,
				mgr.Repos().Settings,
			)
			if err != nil {
				return err
			}

			return jobService.
				NewDailyDecay(
					mgr,
					clock.NewSystemProvider(),
					dailyPeriodResolver,
					app.BuildScoreCalculator(dailyPeriodResolver, cfg.Score),
				).
				ApplyDailyDecay(command.Context())
		})
	},
}

func init() {
	Command.AddCommand(dailyDecayCommand)
}
