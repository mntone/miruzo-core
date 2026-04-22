package database

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func withDropConfirmationIOHook(
	t *testing.T,
	openFn func() (dropConfirmationIO, error),
) {
	t.Helper()

	orig := openDropConfirmationIOFn
	openDropConfirmationIOFn = openFn
	t.Cleanup(func() {
		openDropConfirmationIOFn = orig
	})
}

func TestAskDropConfirmation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "yes_short",
			input: "y\n",
			want:  true,
		},
		{
			name:  "yes_long",
			input: "yes\n",
			want:  true,
		},
		{
			name:  "yes_uppercase",
			input: "YES\n",
			want:  true,
		},
		{
			name:  "no",
			input: "n\n",
			want:  false,
		},
		{
			name:  "empty",
			input: "\n",
			want:  false,
		},
		{
			name:  "eof_without_newline",
			input: "yes",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			promptOutput := &bytes.Buffer{}
			withDropConfirmationIOHook(t, func() (dropConfirmationIO, error) {
				return dropConfirmationIO{
					Reader: strings.NewReader(tt.input),
					Writer: promptOutput,
				}, nil
			})

			got, err := askDropConfirmation()
			assert.NilError(t, "askDropConfirmation() error", err)
			assert.Equal(t, "askDropConfirmation() result", got, tt.want)

			if !strings.Contains(
				promptOutput.String(),
				"Drop database? This operation is destructive. [y/N]: ",
			) {
				t.Fatalf("prompt output = %q", promptOutput.String())
			}
		})
	}
}

func TestAskDropConfirmationReturnsErrorWhenTerminalUnavailable(t *testing.T) {
	terminalErr := errors.New("terminal unavailable")
	withDropConfirmationIOHook(t, func() (dropConfirmationIO, error) {
		return dropConfirmationIO{}, terminalErr
	})

	_, err := askDropConfirmation()
	assert.Error(t, "askDropConfirmation() error", err)
	assert.ErrorIs(
		t,
		"askDropConfirmation() error",
		err,
		errDropConfirmationPromptUnavailable,
	)
	assert.ErrorIs(t, "askDropConfirmation() error", err, terminalErr)
}
