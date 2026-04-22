package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"

	"github.com/charmbracelet/x/term"
)

var (
	ErrStdinNotTerminal = errors.New("stdin is not a terminal")
	ErrEmptyPassword    = errors.New("admin password is empty")
	ErrPasswordEnvUnset = errors.New("admin password environment variable is not set")
	ErrPasswordMode     = errors.New("invalid password input mode")
)

var (
	passwordStdinReader  = io.Reader(os.Stdin)
	passwordStderrWriter = io.Writer(os.Stderr)
	passwordLookupEnv    = func(key string) (string, bool) {
		return os.LookupEnv(key)
	}
	readPasswordTermFn = func() ([]byte, error) {
		return term.ReadPassword(uintptr(syscall.Stdin))
	}
)

func isStdinTerminal() bool {
	return term.IsTerminal(uintptr(syscall.Stdin))
}

func readPasswordFromEnv(envName string) ([]byte, error) {
	if envName == "" {
		return nil, ErrPasswordEnvUnset
	}

	password, ok := passwordLookupEnv(envName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPasswordEnvUnset, envName)
	}

	return []byte(password), nil
}

func readPasswordFromPrompt(prompt string) ([]byte, error) {
	_, err := fmt.Fprint(passwordStderrWriter, prompt)
	if err != nil {
		return nil, fmt.Errorf("write password prompt: %w", err)
	}

	passwordBytes, err := readPasswordTermFn()
	_, newlineErr := fmt.Fprintln(passwordStderrWriter)
	if err != nil {
		return nil, err
	}
	if newlineErr != nil {
		return nil, fmt.Errorf("write password prompt newline: %w", newlineErr)
	}

	return passwordBytes, nil
}

func readPasswordFromStdin() ([]byte, error) {
	passwordBytes, err := io.ReadAll(passwordStdinReader)
	if err != nil {
		return nil, err
	}

	return passwordBytes, nil
}

func normalizePasswordFromStdin(passwordBytes []byte) []byte {
	passwordBytes = bytes.TrimSuffix(passwordBytes, []byte{'\n'})
	passwordBytes = bytes.TrimSuffix(passwordBytes, []byte{'\r'})
	return passwordBytes
}

type ReadPasswordMode int8

const (
	ReadPasswordModeDefault ReadPasswordMode = iota
	ReadPasswordModeStdin
	ReadPasswordModeEnv
)

func (m ReadPasswordMode) String() string {
	switch m {
	case ReadPasswordModeDefault:
		return "default"
	case ReadPasswordModeStdin:
		return "stdin"
	case ReadPasswordModeEnv:
		return "env"
	default:
		return strconv.Itoa(int(m))
	}
}

var (
	isStdinTerminalFn      = isStdinTerminal
	readPasswordFromTTYFn  = readPasswordFromPrompt
	readPasswordFromEnvFn  = readPasswordFromEnv
	readPasswordFromPipeFn = readPasswordFromStdin
)

func ReadPassword(mode ReadPasswordMode, envName string) (string, error) {
	var passwordBytes []byte
	var err error
	switch mode {
	case ReadPasswordModeDefault:
		if !isStdinTerminalFn() {
			return "", ErrStdinNotTerminal
		}

		passwordBytes, err = readPasswordFromTTYFn("Admin password: ")
	case ReadPasswordModeEnv:
		passwordBytes, err = readPasswordFromEnvFn(envName)
	case ReadPasswordModeStdin:
		passwordBytes, err = readPasswordFromPipeFn()
		passwordBytes = normalizePasswordFromStdin(passwordBytes)
	default:
		return "", fmt.Errorf("%w: %s", ErrPasswordMode, mode.String())
	}
	if err != nil {
		return "", err
	}

	if len(passwordBytes) == 0 {
		return "", ErrEmptyPassword
	}

	return string(passwordBytes), nil
}
