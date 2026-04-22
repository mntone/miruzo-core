package database

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	"github.com/spf13/cobra"
)

var errDropConfirmationPromptUnavailable = errors.New(
	"drop confirmation prompt requires a terminal; use --yes to skip confirmation",
)

var dropSkipConfirmation bool

func askDropConfirmation() (bool, error) {
	promptIO, err := openDropConfirmationIOFn()
	if err != nil {
		return false, fmt.Errorf("%w: %w", errDropConfirmationPromptUnavailable, err)
	}
	defer func() {
		_ = promptIO.Close()
	}()

	_, err = fmt.Fprint(
		promptIO.Writer,
		"Drop database? This operation is destructive. [y/N]: ",
	)
	if err != nil {
		return false, err
	}

	line, err := bufio.NewReader(promptIO.Reader).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

var dropCommand = &cobra.Command{
	Use:   "drop",
	Short: "Drop application database",
	Args:  cobra.NoArgs,
	RunE: func(command *cobra.Command, args []string) error {
		return withDatabaseAdminHandle(
			command,
			func(hdl persistence.DatabaseAdminHandle) error {
				exists, err := hdl.Exists(command.Context())
				if err != nil {
					return fmt.Errorf("check database exists: %w", err)
				}
				if !exists {
					return fmt.Errorf("database does not exist")
				}

				if !dropSkipConfirmation {
					confirmed, err := askDropConfirmation()
					if err != nil {
						return err
					}
					if !confirmed {
						_, _ = fmt.Fprintln(command.ErrOrStderr(), "drop canceled")
						return nil
					}
				}

				if err := hdl.Drop(command.Context()); err != nil {
					return err
				}
				return nil
			},
		)
	},
}

func init() {
	dropCommand.Flags().BoolVarP(
		&dropSkipConfirmation,
		"yes",
		"y",
		false,
		"Skip confirmation prompt",
	)
	Command.AddCommand(dropCommand)
}
