package database

import (
	"errors"
	"fmt"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	adaptershared "github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence/shared"
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	cliio "github.com/mntone/miruzo-core/miruzo/internal/cli/io"
	"github.com/mntone/miruzo-core/miruzo/internal/config"
	"github.com/mntone/miruzo-core/miruzo/internal/database/backend"
	"github.com/spf13/cobra"
)

type databaseCommandCallback func(hdl persistence.DatabaseAdminHandle) error

var (
	adminDatabase      string
	adminUserName      string
	adminPassword      string
	adminPasswordEnv   string
	adminPasswordStdin bool
)

var readAdminPasswordFn = cliio.ReadPassword

func validateSQLiteAdminOverrides(cfg config.DatabaseConfig) error {
	if cfg.Backend != backend.SQLite {
		return nil
	}
	if adminDatabase != "" {
		return errors.New("sqlite backend does not support --admin-database")
	}
	if adminUserName != "" {
		return errors.New("sqlite backend does not support --admin-username")
	}
	if cfg.AdminDatabase != "" {
		return errors.New("sqlite backend does not support database.admin_database")
	}
	if cfg.AdminUserName != "" {
		return errors.New("sqlite backend does not support database.admin_username")
	}
	return nil
}

func resolveAdminDatabase(cfg config.DatabaseConfig) string {
	if adminDatabase != "" {
		return adminDatabase
	}
	if cfg.AdminDatabase != "" {
		return cfg.AdminDatabase
	}

	switch cfg.Backend {
	case backend.MySQL:
		return "mysql"
	case backend.PostgreSQL:
		return "postgres"
	}
	return ""
}

func resolveAdminUserName(cfg config.DatabaseConfig) string {
	if adminUserName != "" {
		return adminUserName
	}
	return cfg.AdminUserName
}

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
	if err := validateSQLiteAdminOverrides(cfg.Database); err != nil {
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
			Database: resolveAdminDatabase(cfg.Database),
			UserName: resolveAdminUserName(cfg.Database),
			Password: resolvedAdminPassword,
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
		&adminDatabase,
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
