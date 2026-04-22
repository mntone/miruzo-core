package database

import (
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	adaptershared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	cliio "github.com/mntone/miruzo-core/miruzo/internal/cli/io"
	"github.com/spf13/cobra"
)

type databaseCommandCallback func(hdl persistence.DatabaseAdminHandle) error

var (
	adminDatabaseName  string
	adminUserName      string
	adminPassword      string
	adminPasswordEnv   string
	adminPasswordStdin bool
)

var readAdminPasswordFn = cliio.ReadPassword

func resolveAdminPassword() (string, error) {
	sourceCount := 0
	if adminPassword != "" {
		sourceCount++
	}
	if adminPasswordStdin {
		sourceCount++
	}
	if adminPasswordEnv != "" {
		sourceCount++
	}
	if sourceCount > 1 {
		return "", errors.New(
			"--admin-password, --admin-password-stdin and --admin-password-env are mutually exclusive",
		)
	}

	if adminPasswordStdin {
		password, err := readAdminPasswordFn(
			cliio.ReadPasswordModeStdin,
			"",
		)
		if err != nil && errors.Is(err, cliio.ErrStdinNotTerminal) {
			return "", fmt.Errorf("%w; use --admin-password-stdin", err)
		}
		return password, err
	}

	if adminPasswordEnv != "" {
		password, err := readAdminPasswordFn(
			cliio.ReadPasswordModeEnv,
			adminPasswordEnv,
		)
		return password, err
	}

	return adminPassword, nil
}

func withDatabaseAdminHandle(
	command *cobra.Command,
	callback databaseCommandCallback,
) (err error) {
	cfg, err := app.LoadConfig()
	if err != nil {
		return err
	}

	resolvedAdminPassword, err := resolveAdminPassword()
	if err != nil {
		return err
	}

	hdl, err := persistence.OpenAdminHandle(
		command.Context(),
		cfg.Database,
		adaptershared.DatabaseAdminOptions{
			DatabaseName: adminDatabaseName,
			UserName:     adminUserName,
			Password:     resolvedAdminPassword,
		},
	)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, hdl.Close())
	}()

	return callback(hdl)
}

var Command = &cobra.Command{
	Use:   "database",
	Short: "Database administration commands",
}

func init() {
	Command.PersistentFlags().StringVar(
		&adminDatabaseName,
		"admin-database",
		"",
		"Admin database name used for create/drop operations",
	)
	Command.PersistentFlags().StringVar(
		&adminUserName,
		"admin-username",
		"",
		"Admin username override for create/drop operations",
	)
	Command.PersistentFlags().StringVar(
		&adminPassword,
		"admin-password",
		"",
		"Admin password override for create/drop operations",
	)
	Command.PersistentFlags().StringVar(
		&adminPasswordEnv,
		"admin-password-env",
		"",
		"Environment variable name for admin password",
	)
	Command.PersistentFlags().BoolVar(
		&adminPasswordStdin,
		"admin-password-stdin",
		false,
		"Read admin password from stdin",
	)
}
