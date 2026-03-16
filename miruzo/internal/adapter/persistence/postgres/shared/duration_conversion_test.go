package shared

import (
	"testing"
	"time"
)

func TestPgtypeIntervalFromDurationReturnsValidInterval(t *testing.T) {
	want := 90*time.Second + 123*time.Microsecond

	got := PgtypeIntervalFromDuration(want)

	if !got.Valid {
		t.Fatal("PgtypeIntervalFromDuration(valid) .Valid = false, want true")
	}
	if got.Microseconds != want.Microseconds() {
		t.Fatalf(
			"PgtypeIntervalFromDuration(valid) .Microseconds = %d, want %d",
			got.Microseconds,
			want.Microseconds(),
		)
	}
}
