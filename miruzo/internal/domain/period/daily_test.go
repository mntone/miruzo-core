package period_test

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestDailyResolverPeriodStartKeepsTimeBeforeReset(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "4:59:59.999999999",
			time: time.Date(2026, 1, 2, 4, 59, 59, 999999999, time.UTC),
		},
		{
			name: "4:59:59",
			time: time.Date(2026, 1, 2, 4, 59, 59, 0, time.UTC),
		},
		{
			name: "4:59:00",
			time: time.Date(2026, 1, 2, 4, 59, 0, 0, time.UTC),
		},
		{
			name: "4:00:00",
			time: time.Date(2026, 1, 2, 4, 0, 0, 0, time.UTC),
		},
		{
			name: "23:59:59",
			time: time.Date(2026, 1, 1, 23, 59, 59, 0, time.UTC),
		},
	}

	want := time.Date(2026, 1, 1, 5, 0, 0, 0, time.UTC)
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.PeriodStart(tt.time)
			assert.EqualFn(t, "PeriodStart()", result, want)
		})
	}
}

func TestDailyResolverPeriodStartKeepsTimeAfterReset(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "5:00:00",
			time: time.Date(2026, 1, 2, 5, 0, 0, 0, time.UTC),
		},
		{
			name: "5:01:00",
			time: time.Date(2026, 1, 2, 5, 1, 0, 0, time.UTC),
		},
		{
			name: "5:00:01",
			time: time.Date(2026, 1, 2, 5, 0, 1, 0, time.UTC),
		},
		{
			name: "6:00:00",
			time: time.Date(2026, 1, 2, 6, 0, 0, 0, time.UTC),
		},
		{
			name: "0:00:00",
			time: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	want := time.Date(2026, 1, 2, 5, 0, 0, 0, time.UTC)
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.PeriodStart(tt.time)
			assert.EqualFn(t, "PeriodStart()", result, want)
		})
	}
}

func TestDailyResolverPeriodStartHandlesDSTTransition(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	assert.NilError(t, "time.LoadLocation()", err)

	resolver := period.NewDailyResolverWithLocation(5*time.Hour, location)
	got := resolver.PeriodStart(time.Date(2026, 1, 2, 5, 0, 0, 0, location))
	assert.EqualFn(t, "PeriodStart()", got, time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC))
}

func TestDailyResolverPeriodEnd(t *testing.T) {
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)
	got := resolver.PeriodEnd(time.Date(2026, 1, 2, 6, 0, 0, 0, time.UTC))
	assert.EqualFn(t, "PeriodEnd()", got, time.Date(2026, 1, 3, 5, 0, 0, 0, time.UTC))
}

func TestDailyResolverPeriodRangeReturnsOneDaySpan(t *testing.T) {
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)
	gotStart, gotEnd := resolver.PeriodRange(time.Date(2026, 1, 2, 6, 0, 0, 0, time.UTC))
	assert.EqualFn(t, "PeriodEnd()[0]", gotStart, time.Date(2026, 1, 2, 5, 0, 0, 0, time.UTC))
	assert.EqualFn(t, "PeriodEnd()[1]", gotEnd, time.Date(2026, 1, 3, 5, 0, 0, 0, time.UTC))
}

func TestDailyResolverSincePeriodStartChecksPeriodBoundary(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		flag bool
	}{
		{
			name: "1/1 23:59:59",
			time: time.Date(2026, 1, 1, 23, 59, 59, 0, time.UTC),
			flag: false,
		},
		{
			name: "1/2 4:59:59.999999999",
			time: time.Date(2026, 1, 2, 4, 59, 59, 999999999, time.UTC),
			flag: false,
		},
		{
			name: "1/2 5:00:00",
			time: time.Date(2026, 1, 2, 5, 0, 0, 0, time.UTC),
			flag: true,
		},
		{
			name: "1/2 5:01:00",
			time: time.Date(2026, 1, 2, 5, 1, 0, 0, time.UTC),
			flag: true,
		},
		{
			name: "1/3 0:00:00",
			time: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
			flag: true,
		},
	}

	evaluatedAt := time.Date(2026, 1, 2, 6, 0, 0, 0, time.UTC)
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.SincePeriodStart(tt.time, evaluatedAt)
			assert.Equal(t, "SincePeriodStart()", result, tt.flag)
		})
	}
}

func TestDailyResolverInPeriodChecksPeriodBoundary(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		flag bool
	}{
		{
			name: "1/1 23:59:59",
			time: time.Date(2026, 1, 1, 23, 59, 59, 0, time.UTC),
			flag: false,
		},
		{
			name: "1/2 4:59:59.999999999",
			time: time.Date(2026, 1, 2, 4, 59, 59, 999999999, time.UTC),
			flag: false,
		},
		{
			name: "1/2 5:00:00",
			time: time.Date(2026, 1, 2, 5, 0, 0, 0, time.UTC),
			flag: true,
		},
		{
			name: "1/2 5:01:00",
			time: time.Date(2026, 1, 2, 5, 1, 0, 0, time.UTC),
			flag: true,
		},
		{
			name: "1/3 0:00:00",
			time: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
			flag: true,
		},
		{
			name: "1/3 4:59:59.999999999",
			time: time.Date(2026, 1, 3, 4, 59, 59, 999999999, time.UTC),
			flag: true,
		},
		{
			name: "1/3 5:00:00",
			time: time.Date(2026, 1, 3, 5, 0, 0, 0, time.UTC),
			flag: false,
		},
	}

	evaluatedAt := time.Date(2026, 1, 2, 6, 0, 0, 0, time.UTC)
	resolver := period.NewDailyResolverWithLocation(5*time.Hour, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.InPeriod(tt.time, evaluatedAt)
			assert.Equal(t, "InPeriod()", result, tt.flag)
		})
	}
}
