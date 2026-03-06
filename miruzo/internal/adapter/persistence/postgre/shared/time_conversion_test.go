package shared

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/mo"
)

func TestTimeFromPgtypeReturnsTimeWhenValid(t *testing.T) {
	want := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	got := TimeFromPgtype(pgtype.Timestamp{
		Time:  want,
		Valid: true,
	})

	if !got.Equal(want) {
		t.Fatalf("TimeFromPgtype(valid) = %s, want %s", got, want)
	}
}

func TestTimeFromPgtypePanicsWhenInvalid(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("TimeFromPgtype(invalid) did not panic")
		}

		recoveredErr, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic value type = %T, want error", recovered)
		}
		if !errors.Is(recoveredErr, errOptionNoSuchElement) {
			t.Fatalf("panic error = %v, want %v", recoveredErr, errOptionNoSuchElement)
		}
	}()

	_ = TimeFromPgtype(pgtype.Timestamp{})
}

func TestPgtypeTimestampFromTimeReturnsValidTimestamp(t *testing.T) {
	want := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	got := PgtypeTimestampFromTime(want)

	if !got.Valid {
		t.Fatal("PgtypeTimestampFromTime(valid) .Valid = false, want true")
	}
	if !got.Time.Equal(want) {
		t.Fatalf("PgtypeTimestampFromTime(valid) .Time = %s, want %s", got.Time, want)
	}
}

func TestPgtypeTimestampFromOptionNoneReturnsInvalidTimestamp(t *testing.T) {
	got := PgtypeTimestampFromOption(mo.None[time.Time]())

	if got.Valid {
		t.Fatal("PgtypeTimestampFromOption(None) .Valid = true, want false")
	}
}

func TestPgtypeTimestampFromOptionSomeReturnsValidTimestamp(t *testing.T) {
	want := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

	got := PgtypeTimestampFromOption(mo.Some(want))

	if !got.Valid {
		t.Fatal("PgtypeTimestampFromOption(Some) .Valid = false, want true")
	}
	if !got.Time.Equal(want) {
		t.Fatalf("PgtypeTimestampFromOption(Some) .Time = %s, want %s", got.Time, want)
	}
}
