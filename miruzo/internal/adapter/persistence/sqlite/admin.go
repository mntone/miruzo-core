package sqlite

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func sqliteDatabaseFilePath(dsn string) (string, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	var rawPath string
	switch {
	case parsed.Scheme == "file":
		if parsed.Opaque != "" {
			rawPath = parsed.Opaque
		} else {
			rawPath = parsed.Path
		}
	case parsed.Scheme == "":
		rawPath = parsed.Path
	default:
		return "", fmt.Errorf("unsupported sqlite dsn scheme: %s", parsed.Scheme)
	}

	path, err := url.PathUnescape(rawPath)
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", errors.New("sqlite dsn path must not be empty")
	}
	if path == ":memory:" || parsed.Query().Get("mode") == "memory" {
		return "", errors.New("sqlite memory database is not supported for DatabaseAdmin")
	}

	return filepath.Clean(path), nil
}

type sqliteAdminHandle struct {
	filePath string
}

func OpenAdminHandle(
	appConfig config.DatabaseConfig,
	adminDatabaseName string,
) (sqliteAdminHandle, error) {
	if adminDatabaseName == "" {
		adminDatabaseName = appConfig.AdminDatabaseName
	}
	if adminDatabaseName != "" {
		return sqliteAdminHandle{}, fmt.Errorf(
			"sqlite backend does not support admin database override: %q",
			adminDatabaseName,
		)
	}

	filePath, err := sqliteDatabaseFilePath(appConfig.DSN)
	if err != nil {
		return sqliteAdminHandle{}, err
	}

	return sqliteAdminHandle{
		filePath: filePath,
	}, nil
}

func (hdl sqliteAdminHandle) Close() error {
	return nil
}

func (hdl sqliteAdminHandle) Create(_ context.Context) error {
	file, err := os.OpenFile(hdl.filePath, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf(
			"sqlite admin create database %q failed: %w",
			hdl.filePath,
			err,
		)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf(
			"sqlite admin create database %q failed: %w",
			hdl.filePath,
			err,
		)
	}
	return nil
}

func (hdl sqliteAdminHandle) Drop(_ context.Context) error {
	if err := os.Remove(hdl.filePath); err != nil {
		return fmt.Errorf(
			"sqlite admin drop database %q failed: %w",
			hdl.filePath,
			err,
		)
	}
	return nil
}

func (hdl sqliteAdminHandle) Exists(_ context.Context) (bool, error) {
	_, err := os.Stat(hdl.filePath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf(
		"sqlite admin check database %q exists failed: %w",
		hdl.filePath,
		err,
	)
}
