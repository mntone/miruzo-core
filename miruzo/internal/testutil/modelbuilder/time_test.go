package modelbuilder

import (
	"testing"
	"time"

	"github.com/samber/mo"
)

func TestResolveOffsetTime(t *testing.T) {
	baseTime := time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC)

	t.Run("time.Duration", func(t *testing.T) {
		got, present := resolveOffsetTime(2*time.Minute, baseTime).Get()
		if !present {
			t.Fatal("resolveOffsetTime() = none, want some")
		}
		want := baseTime.Add(2 * time.Minute)
		if got != want {
			t.Fatalf("resolveOffsetTime() = %v, want %v", got, want)
		}
	})

	integerTests := []struct {
		name  string
		value any
	}{
		{
			name:  "int seconds",
			value: 30,
		},
		{
			name:  "int8 seconds",
			value: int8(30),
		},
		{
			name:  "int16 seconds",
			value: int16(30),
		},
		{
			name:  "int32 seconds",
			value: int32(30),
		},
		{
			name:  "int64 seconds",
			value: int64(30),
		},
		{
			name:  "uint seconds",
			value: uint(30),
		},
		{
			name:  "uint8 seconds",
			value: uint8(30),
		},
		{
			name:  "uint16 seconds",
			value: uint16(30),
		},
		{
			name:  "uint32 seconds",
			value: uint32(30),
		},
		{
			name:  "uint64 seconds",
			value: uint64(30),
		},
	}
	for _, tt := range integerTests {
		t.Run(tt.name, func(t *testing.T) {
			got, present := resolveOffsetTime(tt.value, baseTime).Get()
			if !present {
				t.Fatal("resolveOffsetTime() = none, want some")
			}
			want := baseTime.Add(30 * time.Second)
			if got != want {
				t.Fatalf("resolveOffsetTime() = %v, want %v", got, want)
			}
		})
	}

	t.Run("option some", func(t *testing.T) {
		got, present := resolveOffsetTime(mo.Some(45*time.Second), baseTime).Get()
		if !present {
			t.Fatal("resolveOffsetTime() = none, want some")
		}
		want := baseTime.Add(45 * time.Second)
		if got != want {
			t.Fatalf("resolveOffsetTime() = %v, want %v", got, want)
		}
	})

	t.Run("option none", func(t *testing.T) {
		_, present := resolveOffsetTime(mo.None[time.Duration](), baseTime).Get()
		if present {
			t.Fatal("resolveOffsetTime() = some, want none")
		}
	})

	stringTests := []struct {
		name   string
		value  string
		offset time.Duration
	}{
		{
			name:   "nanoseconds",
			value:  "400ns",
			offset: 400 * time.Nanosecond,
		},
		{
			name:   "microseconds alphabet",
			value:  "80us",
			offset: 80 * time.Microsecond,
		},
		{
			name:   "microseconds",
			value:  "90µs",
			offset: 90 * time.Microsecond,
		},
		{
			name:   "milliseconds",
			value:  "200ms",
			offset: 200 * time.Millisecond,
		},
		{
			name:   "seconds",
			value:  "3s",
			offset: 3 * time.Second,
		},
		{
			name:   "minutes",
			value:  "5m",
			offset: 5 * time.Minute,
		},
		{
			name:   "hours",
			value:  "-1h",
			offset: -1 * time.Hour,
		},
		{
			name:   "hours and minutes",
			value:  "1h40m",
			offset: time.Hour + 40*time.Minute,
		},
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			got, present := resolveOffsetTime(tt.value, baseTime).Get()
			if !present {
				t.Fatal("resolveOffsetTime() = none, want some")
			}
			want := baseTime.Add(tt.offset)
			if got != want {
				t.Fatalf("resolveOffsetTime() = %v, want %v", got, want)
			}
		})
	}
}

func TestResolveOffsetTimePanicsForInvalidType(t *testing.T) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("resolveOffsetTime() panic = nil, want non-nil")
		}
	}()

	baseTime := time.Date(2026, 1, 10, 5, 0, 0, 0, time.UTC)
	resolveOffsetTime("invalid", baseTime)
}
