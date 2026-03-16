package settings

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/timezone"
	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/samber/mo"
)

const timezoneSettingsKey = "timezone"

type TimezoneProvider struct {
	repository persist.SettingsRepository
	resolver   timezone.TimezoneResolver
	location   mo.Option[string]
}

func NewTimezoneProvider(
	repository persist.SettingsRepository,
	resolver timezone.TimezoneResolver,
) *TimezoneProvider {
	return &TimezoneProvider{
		repository: repository,
		resolver:   resolver,
	}
}

func (prov *TimezoneProvider) Location() string {
	value, present := prov.location.Get()
	if !present {
		log.Fatalf("settings: timezone is not initialized; call EnsureSettings first")
	}

	return value
}

func (prov *TimezoneProvider) SetLocation(location string) {
	prov.location = mo.Some(location)
}

func (prov *TimezoneProvider) saveLocation(
	ctx context.Context,
	location string,
) {
	err := prov.repository.UpdateValue(ctx, timezoneSettingsKey, location)
	if err != nil {
		log.Printf(
			"settings: failed to save timezone: key=%s value=%s",
			timezoneSettingsKey,
			location,
		)
	}
}

func (prov *TimezoneProvider) EnsureSettings(
	ctx context.Context,
	initialLocation *string,
) {
	// 1. Use the saved location from the settings table.
	appLocation, err := prov.repository.GetValue(ctx, timezoneSettingsKey)
	if err == nil {
		prov.location = mo.Some(appLocation)
		return
	}
	if !errors.Is(err, persist.ErrNotFound) {
		log.Fatalf("settings: failed to load timezone: %v", err)
	}

	// 2. If missing, use/save the initial location.
	if initialLocation != nil {
		prov.saveLocation(ctx, *initialLocation)
		prov.location = mo.Some(*initialLocation)
		return
	}

	// 3. If still missing, use/save the host's location.
	systemLocation, present := prov.resolver.GetLocation().Get()
	if present {
		prov.saveLocation(ctx, systemLocation)
		prov.location = mo.Some(systemLocation)
		return
	}

	// 4. If still missing, use/save UTC.
	utcLocation := time.UTC.String()
	prov.saveLocation(ctx, utcLocation)
	prov.location = mo.Some(utcLocation)
}
