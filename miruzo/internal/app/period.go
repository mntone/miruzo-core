package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/timezone"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/service/settings"
)

const defaultDayStartOffsetLocal = 5 * time.Hour

func newDailyResolver(
	ctx context.Context,
	cfg config.PeriodConfig,
	settingsRepository persist.SettingsRepository,
) (period.DailyResolver, error) {
	if cfg.DayStartOffset != nil {
		return period.NewDailyResolver(*cfg.DayStartOffset), nil
	}

	timezoneProvider := settings.NewTimezoneProvider(
		settingsRepository,
		timezone.NewLocalTimezoneResolver(),
	)
	timezoneProvider.EnsureSettings(ctx, cfg.InitialLocation)

	location, err := time.LoadLocation(timezoneProvider.Location())
	if err != nil {
		return period.DailyResolver{}, fmt.Errorf(
			"load location: %s: %w",
			timezoneProvider.Location(),
			err,
		)
	}

	_, offsetSeconds := time.Date(2026, 1, 1, 0, 0, 0, 0, location).Zone()
	timezoneOffset := time.Duration(offsetSeconds) * time.Second

	dayStartOffset := defaultDayStartOffsetLocal - timezoneOffset
	if dayStartOffset < 0 {
		dayStartOffset += 24 * time.Hour
	}

	return period.NewDailyResolver(dayStartOffset), nil
}
