package shared

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/samber/mo"
)

func TestTimeFromSqlReturnsTimeWhenValid(t *testing.T) {
	want := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	got := TimeFromSql(sql.NullTime{
		Time:  want,
		Valid: true,
	})

	if !got.Equal(want) {
		t.Fatalf("TimeFromSql(valid) = %s, want %s", got, want)
	}
}

func TestTimeFromSqlPanicsWhenInvalid(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("TimeFromSql(invalid) did not panic")
		}

		recoveredErr, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic value type = %T, want error", recovered)
		}
		if !errors.Is(recoveredErr, errOptionNoSuchElement) {
			t.Fatalf("panic error = %v, want %v", recoveredErr, errOptionNoSuchElement)
		}
	}()

	_ = TimeFromSql(sql.NullTime{})
}

func TestNullTimeFromTimeReturnsValidNullTime(t *testing.T) {
	want := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	got := NullTimeFromTime(want)

	if !got.Valid {
		t.Fatal("NullTimeFromTime(valid) .Valid = false, want true")
	}
	if !got.Time.Equal(want) {
		t.Fatalf("NullTimeFromTime(valid) .Time = %s, want %s", got.Time, want)
	}
}

func TestNullTimeFromOptionNoneReturnsInvalidNullTime(t *testing.T) {
	got := NullTimeFromOption(mo.None[time.Time]())

	if got.Valid {
		t.Fatal("NullTimeFromOption(None) .Valid = true, want false")
	}
}

func TestNullTimeFromOptionSomeReturnsValidNullTime(t *testing.T) {
	want := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

	got := NullTimeFromOption(mo.Some(want))

	if !got.Valid {
		t.Fatal("NullTimeFromOption(Some) .Valid = false, want true")
	}
	if !got.Time.Equal(want) {
		t.Fatalf("NullTimeFromOption(Some) .Time = %s, want %s", got.Time, want)
	}
}
