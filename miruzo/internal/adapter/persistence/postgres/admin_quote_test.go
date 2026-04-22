package postgres

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestPostgresQuoteIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		identifier string
		want       string
	}{
		{
			name:       "plain_identifier",
			identifier: "miruzo",
			want:       `"miruzo"`,
		},
		{
			name:       "identifier_with_double_quote",
			identifier: `miru"zo`,
			want:       `"miru""zo"`,
		},
		{
			name:       "empty_identifier",
			identifier: "",
			want:       `""`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := postgresQuoteIdentifier(tt.identifier)
			assert.Equal(t, "postgresQuoteIdentifier()", got, tt.want)
		})
	}
}
