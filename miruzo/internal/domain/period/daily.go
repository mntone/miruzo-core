package period

import "time"

type DailyResolver struct {
	location       *time.Location
	dayStartOffset time.Duration
}

func NewDailyResolverWithLocation(
	dayStartOffset time.Duration,
	location *time.Location,
) DailyResolver {
	if location == nil {
		location = time.Local
	}

	return DailyResolver{
		location:       location,
		dayStartOffset: dayStartOffset,
	}
}

func NewDailyResolver(dayStartOffset time.Duration) DailyResolver {
	return NewDailyResolverWithLocation(dayStartOffset, time.Local)
}

// PeriodStart returns the start time of the daily period that contains evaluatedAt.
// The start is determined by the resolver's dayStartOffset in its location.
func (resolv DailyResolver) PeriodStart(evaluatedAt time.Time) time.Time {
	target := evaluatedAt.In(resolv.location)

	midnight := time.Date(
		target.Year(),
		target.Month(),
		target.Day(),
		0, 0, 0, 0,
		resolv.location,
	)

	candidate := midnight.Add(resolv.dayStartOffset)
	if candidate.After(target) {
		candidate = candidate.AddDate(0, 0, -1)
	}

	return candidate
}

// PeriodEnd returns the exclusive end of the daily period containing evaluatedAt.
func (resolv DailyResolver) PeriodEnd(evaluatedAt time.Time) time.Time {
	return resolv.PeriodStart(evaluatedAt).AddDate(0, 0, 1)
}

// PeriodRange returns the daily period containing evaluatedAt as [start, end).
func (resolv DailyResolver) PeriodRange(evaluatedAt time.Time) (start time.Time, end time.Time) {
	start = resolv.PeriodStart(evaluatedAt)
	end = start.AddDate(0, 0, 1)
	return
}

// SincePeriodStart reports whether target is at or after the start of the
// daily period that contains evaluatedAt.
func (resolv DailyResolver) SincePeriodStart(
	target time.Time,
	evaluatedAt time.Time,
) bool {
	start := resolv.PeriodStart(evaluatedAt)
	return !target.Before(start)
}

// InPeriod reports whether target is within the daily period containing evaluatedAt.
func (resolv DailyResolver) InPeriod(
	target time.Time,
	evaluatedAt time.Time,
) bool {
	start := resolv.PeriodStart(evaluatedAt)
	end := start.AddDate(0, 0, 1)
	return !target.Before(start) && target.Before(end)
}
