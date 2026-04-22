package io

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type errorWriter struct {
	err error
}

func (w errorWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type writeErrorOnNth struct {
	n      int
	count  int
	writer io.Writer
	err    error
}

func (w *writeErrorOnNth) Write(p []byte) (int, error) {
	w.count++
	if w.count == w.n {
		return 0, w.err
	}
	return w.writer.Write(p)
}

func withPasswordHooks(t *testing.T) {
	t.Helper()

	origStdin := passwordStdinReader
	origStderr := passwordStderrWriter
	origLookupEnv := passwordLookupEnv
	origReadTerm := readPasswordTermFn
	origIsTerm := isStdinTerminalFn
	origReadTTY := readPasswordFromTTYFn
	origReadPipe := readPasswordFromPipeFn
	origReadEnv := readPasswordFromEnvFn

	t.Cleanup(func() {
		passwordStdinReader = origStdin
		passwordStderrWriter = origStderr
		passwordLookupEnv = origLookupEnv
		readPasswordTermFn = origReadTerm
		isStdinTerminalFn = origIsTerm
		readPasswordFromTTYFn = origReadTTY
		readPasswordFromPipeFn = origReadPipe
		readPasswordFromEnvFn = origReadEnv
	})
}

func TestReadPasswordFromPromptWriteError(t *testing.T) {
	withPasswordHooks(t)

	wantErr := errors.New("stderr write failed")
	passwordStderrWriter = errorWriter{err: wantErr}
	readPasswordTermFn = func() ([]byte, error) {
		return []byte("secret"), nil
	}

	_, err := readPasswordFromPrompt("Admin password: ")
	assert.ErrorIs(t, "readPasswordFromPrompt() error", err, wantErr)
}

func TestReadPasswordFromPromptNewlineWriteError(t *testing.T) {
	withPasswordHooks(t)

	wantErr := errors.New("stderr newline write failed")
	passwordStderrWriter = &writeErrorOnNth{
		n:      2,
		writer: io.Discard,
		err:    wantErr,
	}
	readPasswordTermFn = func() ([]byte, error) {
		return []byte("secret"), nil
	}

	_, err := readPasswordFromPrompt("Admin password: ")
	assert.ErrorIs(t, "readPasswordFromPrompt() error", err, wantErr)
}

func TestReadPasswordFromStdinTrimsTrailingLineBreakOnly(t *testing.T) {
	withPasswordHooks(t)

	readPasswordFromPipeFn = func() ([]byte, error) {
		return []byte("  secret with space  \n"), nil
	}

	got, err := ReadPassword(ReadPasswordModeStdin, "")
	assert.NilError(t, "ReadPassword(stdin) error", err)
	assert.Equal(t, "ReadPassword(stdin)", got, "  secret with space  ")
}

func TestReadPasswordFromStdinTrimsCRLF(t *testing.T) {
	withPasswordHooks(t)

	readPasswordFromPipeFn = func() ([]byte, error) {
		return []byte("secret\r\n"), nil
	}

	got, err := ReadPassword(ReadPasswordModeStdin, "")
	assert.NilError(t, "ReadPassword(stdin) error", err)
	assert.Equal(t, "ReadPassword(stdin)", got, "secret")
}

func TestReadPasswordFromStdinEmpty(t *testing.T) {
	withPasswordHooks(t)

	readPasswordFromPipeFn = func() ([]byte, error) {
		return []byte("\n"), nil
	}

	_, err := ReadPassword(ReadPasswordModeStdin, "")
	assert.Error(t, "ReadPassword(stdin) error", err)
	assert.ErrorIs(t, "ReadPassword(stdin) error", err, ErrEmptyPassword)
}

func TestReadPasswordRequiresTerminalInDefaultMode(t *testing.T) {
	withPasswordHooks(t)

	isStdinTerminalFn = func() bool {
		return false
	}

	_, err := ReadPassword(ReadPasswordModeDefault, "")
	assert.Error(t, "ReadPassword(default) error", err)
	assert.ErrorIs(t, "ReadPassword(default) error", err, ErrStdinNotTerminal)
}

func TestReadPasswordFromPromptPreservesSpaces(t *testing.T) {
	withPasswordHooks(t)

	isStdinTerminalFn = func() bool {
		return true
	}
	readPasswordFromTTYFn = func(_ string) ([]byte, error) {
		return []byte("  secret  "), nil
	}

	got, err := ReadPassword(ReadPasswordModeDefault, "")
	assert.NilError(t, "ReadPassword(default) error", err)
	assert.Equal(t, "ReadPassword(default)", got, "  secret  ")
}

func TestReadPasswordFromEnv(t *testing.T) {
	withPasswordHooks(t)

	passwordLookupEnv = func(key string) (string, bool) {
		assert.Equal(t, "env key", key, "MIRUZO_ADMIN_PASSWORD")
		return "env-secret", true
	}

	got, err := ReadPassword(ReadPasswordModeEnv, "MIRUZO_ADMIN_PASSWORD")
	assert.NilError(t, "ReadPassword(env) error", err)
	assert.Equal(t, "ReadPassword(env)", got, "env-secret")
}

func TestReadPasswordFromEnvUnset(t *testing.T) {
	withPasswordHooks(t)

	passwordLookupEnv = func(key string) (string, bool) {
		return "", false
	}

	_, err := ReadPassword(ReadPasswordModeEnv, "MIRUZO_ADMIN_PASSWORD")
	assert.Error(t, "ReadPassword(env) error", err)
	assert.ErrorIs(t, "ReadPassword(env) error", err, ErrPasswordEnvUnset)
}

func TestReadPasswordFromEnvEmptyName(t *testing.T) {
	withPasswordHooks(t)

	_, err := ReadPassword(ReadPasswordModeEnv, "")
	assert.Error(t, "ReadPassword(env) error", err)
	assert.ErrorIs(t, "ReadPassword(env) error", err, ErrPasswordEnvUnset)
}

func TestReadPasswordInvalidMode(t *testing.T) {
	withPasswordHooks(t)

	_, err := ReadPassword(ReadPasswordMode(-1), "")
	assert.Error(t, "ReadPassword(unknown) error", err)
	assert.ErrorIs(t, "ReadPassword(unknown) error", err, ErrPasswordMode)
}

func TestNormalizePasswordFromStdin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "no_linebreak", input: "secret", want: "secret"},
		{name: "lf", input: "secret\n", want: "secret"},
		{name: "crlf", input: "secret\r\n", want: "secret"},
		{name: "preserve_trailing_spaces", input: "secret  \n", want: "secret  "},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePasswordFromStdin([]byte(tt.input))
			if strings.Compare(string(got), tt.want) != 0 {
				t.Fatalf("normalizePasswordFromStdin() = %q, want %q", string(got), tt.want)
			}
		})
	}
}
