package score_test

import (
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/domain/period"
	"github.com/mntone/miruzo-core/miruzo/internal/domain/score"
	"github.com/mntone/miruzo-core/miruzo/internal/model"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
	"github.com/samber/mo"
)

var calc = score.New(
	period.NewDailyResolverWithLocation(
		5*time.Hour,
		time.UTC,
	),
	10,
	[]model.ScoreViewBonusRule{
		{Days: 1, Bonus: 10},
		{Days: 3, Bonus: 7},
		{Days: 7, Bonus: 5},
		{Days: 30, Bonus: 3},
	},
	2,
	1, -1,
	20, -18,
)

func TestScoreCalculatorViewDelta(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want model.ScoreType
	}{
		{
			name: "Days=1",
			time: time.Date(2026, 1, 29, 3, 4, 30, 0, time.UTC),
			want: 10,
		},
		{
			name: "Days=2",
			time: time.Date(2026, 1, 28, 3, 4, 30, 0, time.UTC),
			want: 7,
		},
		{
			name: "Days=3",
			time: time.Date(2026, 1, 27, 3, 4, 30, 0, time.UTC),
			want: 7,
		},
		{
			name: "Days=5",
			time: time.Date(2026, 1, 25, 3, 4, 30, 0, time.UTC),
			want: 5,
		},
		{
			name: "Days=7",
			time: time.Date(2026, 1, 23, 3, 4, 30, 0, time.UTC),
			want: 5,
		},
		{
			name: "Days=8",
			time: time.Date(2026, 1, 22, 3, 4, 30, 0, time.UTC),
			want: 3,
		},
		{
			name: "Days=29",
			time: time.Date(2026, 1, 1, 3, 4, 30, 0, time.UTC),
			want: 3,
		},
		{
			name: "Days=30",
			time: time.Date(2025, 12, 31, 3, 4, 30, 0, time.UTC),
			want: 3,
		},
		{
			name: "Days=31",
			time: time.Date(2025, 12, 30, 3, 4, 30, 0, time.UTC),
			want: 2,
		},
	}

	evaluatedAt := time.Date(2026, 1, 29, 5, 0, 0, 0, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, negative := calc.ViewDelta(mo.Some(tt.time), evaluatedAt)
			assert.Equal(t, "ViewDelta()[0]", got, tt.want)
			assert.Equal(t, "ViewDelta()[1]", negative, false)
		})
	}
}

func TestScoreCalculatorViewDeltaGetsFirstBonus(t *testing.T) {
	delta, negative := calc.ViewDelta(
		mo.None[time.Time](),
		time.Date(2026, 1, 30, 4, 5, 0, 0, time.UTC),
	)

	assert.Equal(t, "ViewDelta()[0]", delta, 10)
	assert.Equal(t, "ViewDelta()[1]", negative, false)
}

func TestScoreCalculatorViewDeltaGetsNegative(t *testing.T) {
	delta, negative := calc.ViewDelta(
		mo.Some(time.Date(2026, 1, 31, 4, 4, 30, 0, time.UTC)),
		time.Date(2026, 1, 30, 4, 5, 0, 0, time.UTC),
	)

	assert.Equal(t, "ViewDelta()[0]", delta, 0)
	assert.Equal(t, "ViewDelta()[1]", negative, true)
}
