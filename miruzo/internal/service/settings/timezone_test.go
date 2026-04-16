package settings

import (
	"context"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/persist"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/persistence/stub"
	tt "github.com/mntone/miruzo-core/miruzo/internal/testutil/adapter/timezone"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

const timezoneKey = "timezone"

func TestTimezoneProviderEnsureSettingsUsesStoredValue(t *testing.T) {
	t.Parallel()

	repo := stub.NewStubSettingsRepositoryWithKeyValue(timezoneKey, "Asia/Tokyo")
	resolver := tt.NewFixedTimezoneResolverWithLocation("America/New_York")
	provider := NewTimezoneProvider(repo, resolver)

	provider.EnsureSettings(context.Background(), nil)

	assert.Equal(t, "Location()", provider.Location(), "Asia/Tokyo")
	assert.Empty(t, "Updates", repo.Updates)
}

func TestTimezoneProviderEnsureSettingsUsesInitialLocation(t *testing.T) {
	t.Parallel()

	repo := stub.NewStubSettingsRepositoryWithGetError(persist.ErrNoRows)
	resolver := tt.NewFixedTimezoneResolverWithLocation("America/New_York")
	provider := NewTimezoneProvider(repo, resolver)

	initial := "Asia/Tokyo"
	provider.EnsureSettings(context.Background(), &initial)

	assert.Equal(t, "Location()", provider.Location(), "Asia/Tokyo")
	assert.LenIs(t, "Updates", repo.Updates, 1)
	assert.Equal(t, "Update[0].Value", repo.Updates[0].Value, "Asia/Tokyo")
}

func TestTimezoneProviderEnsureSettingsUsesSystemLocation(t *testing.T) {
	t.Parallel()

	repo := stub.NewStubSettingsRepositoryWithGetError(persist.ErrNoRows)
	resolver := tt.NewFixedTimezoneResolverWithLocation("America/New_York")
	provider := NewTimezoneProvider(repo, resolver)

	provider.EnsureSettings(context.Background(), nil)

	assert.Equal(t, "Location()", provider.Location(), "America/New_York")
	assert.LenIs(t, "Updates", repo.Updates, 1)
	assert.Equal(t, "Update[0].Value", repo.Updates[0].Value, "America/New_York")
}

func TestTimezoneProviderEnsureSettingsFallsBackToUTC(t *testing.T) {
	t.Parallel()

	repo := stub.NewStubSettingsRepositoryWithGetError(persist.ErrNoRows)
	resolver := tt.NewEmptyTimezoneResolver()
	provider := NewTimezoneProvider(repo, resolver)

	provider.EnsureSettings(context.Background(), nil)

	assert.Equal(t, "Location()", provider.Location(), time.UTC.String())
	assert.LenIs(t, "Updates", repo.Updates, 1)
	assert.Equal(t, "Update[0].Value", repo.Updates[0].Value, time.UTC.String())
}
