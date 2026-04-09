package shared

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/mntone/miruzo-core/miruzo/internal/database/migration"
)

type MigrationRunner struct {
	migration.Spec
}

func (r MigrationRunner) Migrate(ctx context.Context, version int) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	return m.Migrate(uint(version))
}

func (r MigrationRunner) Step(ctx context.Context, steps int) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	err = m.Steps(steps)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r MigrationRunner) Down(ctx context.Context) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r MigrationRunner) Up(ctx context.Context) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r MigrationRunner) Version(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	version, _, err := m.Version()
	if err != nil {
		return 0, err
	}

	return int(version), nil
}

func (r MigrationRunner) SetVersion(ctx context.Context, version int) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	m, close, err := r.NewInstance()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, close())
	}()

	return m.Force(version)
}
