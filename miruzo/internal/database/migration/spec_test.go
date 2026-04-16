package migration

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

type fakeSourceDriver struct {
	closeCount int
	closeErr   error
}

func (driver *fakeSourceDriver) Open(_ string) (source.Driver, error) {
	return driver, nil
}

func (driver *fakeSourceDriver) Close() error {
	driver.closeCount++
	return driver.closeErr
}

func (driver *fakeSourceDriver) First() (uint, error) {
	return 0, os.ErrNotExist
}

func (driver *fakeSourceDriver) Prev(uint) (uint, error) {
	return 0, os.ErrNotExist
}

func (driver *fakeSourceDriver) Next(uint) (uint, error) {
	return 0, os.ErrNotExist
}

func (driver *fakeSourceDriver) ReadUp(uint) (io.ReadCloser, string, error) {
	return nil, "", os.ErrNotExist
}

func (driver *fakeSourceDriver) ReadDown(uint) (io.ReadCloser, string, error) {
	return nil, "", os.ErrNotExist
}

type fakeDatabaseDriver struct {
	closeCount int
	closeErr   error
}

func (driver *fakeDatabaseDriver) Open(_ string) (database.Driver, error) {
	return driver, nil
}

func (driver *fakeDatabaseDriver) Close() error {
	driver.closeCount++
	return driver.closeErr
}

func (driver *fakeDatabaseDriver) Lock() error {
	return nil
}

func (driver *fakeDatabaseDriver) Unlock() error {
	return nil
}

func (driver *fakeDatabaseDriver) Run(io.Reader) error {
	return nil
}

func (driver *fakeDatabaseDriver) SetVersion(int, bool) error {
	return nil
}

func (driver *fakeDatabaseDriver) Version() (int, bool, error) {
	return -1, false, nil
}

func (driver *fakeDatabaseDriver) Drop() error {
	return nil
}

func TestSpecNewInstanceReturnsErrorWhenSourceCreateFails(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("source failed")
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return nil, expectedErr
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			t.Fatalf("database driver factory should not be called")
			return nil, nil
		},
	}

	m, close, err := spec.NewInstance()
	assert.Error(t, "NewInstance() error", err)
	assert.Nil(t, "NewInstance()", m)
	if close != nil {
		t.Fatalf("expected nil close func, got non-nil")
	}
	if !strings.Contains(err.Error(), "create migration source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpecNewInstanceClosesSourceWhenDatabaseCreateFails(t *testing.T) {
	t.Parallel()

	sourceDriver := &fakeSourceDriver{}
	expectedErr := errors.New("database failed")
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return sourceDriver, nil
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			return nil, expectedErr
		},
	}

	m, close, err := spec.NewInstance()
	assert.Error(t, "NewInstance() error", err)
	assert.Nil(t, "NewInstance()", m)
	if close != nil {
		t.Fatalf("expected nil close func, got non-nil")
	}
	assert.Equal(t, "sourceDriver.closeCount", sourceDriver.closeCount, 1)
	if !strings.Contains(err.Error(), "create migration driver") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpecNewInstanceCloseOnlySourceWhenKeepDatabaseOpenTrue(t *testing.T) {
	t.Parallel()

	sourceDriver := &fakeSourceDriver{}
	databaseDriver := &fakeDatabaseDriver{}
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return sourceDriver, nil
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			return databaseDriver, nil
		},
		KeepDatabaseOpen: true,
	}

	m, close, err := spec.NewInstance()
	assert.NilError(t, "NewInstance() error", err)
	assert.NotNil(t, "NewInstance()", m)
	if close == nil {
		t.Fatalf("expected close func, got nil")
	}

	closeErr := close()
	assert.NilError(t, "close() error", closeErr)
	assert.Equal(t, "sourceDriver.closeCount", sourceDriver.closeCount, 1)
	assert.Equal(t, "databaseDriver.closeCount", databaseDriver.closeCount, 0)
}

func TestSpecNewInstanceCloseSourceAndDatabaseWhenKeepDatabaseOpenFalse(t *testing.T) {
	t.Parallel()

	sourceDriver := &fakeSourceDriver{}
	databaseDriver := &fakeDatabaseDriver{}
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return sourceDriver, nil
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			return databaseDriver, nil
		},
		KeepDatabaseOpen: false,
	}

	m, close, err := spec.NewInstance()
	assert.NilError(t, "NewInstance() error", err)
	assert.NotNil(t, "NewInstance()", m)
	if close == nil {
		t.Fatalf("expected close func, got nil")
	}

	closeErr := close()
	assert.NilError(t, "close() error", closeErr)
	assert.Equal(t, "sourceDriver.closeCount", sourceDriver.closeCount, 1)
	assert.Equal(t, "databaseDriver.closeCount", databaseDriver.closeCount, 1)
}

func TestSpecNewInstanceCloseReturnsSourceErrorWhenKeepDatabaseOpenTrue(t *testing.T) {
	t.Parallel()

	sourceDriver := &fakeSourceDriver{
		closeErr: errors.New("source close failed"),
	}
	databaseDriver := &fakeDatabaseDriver{}
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return sourceDriver, nil
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			return databaseDriver, nil
		},
		KeepDatabaseOpen: true,
	}

	m, close, err := spec.NewInstance()
	assert.NilError(t, "NewInstance() error", err)
	assert.NotNil(t, "NewInstance()", m)
	if close == nil {
		t.Fatalf("close() = nil, want non-nil")
	}

	closeErr := close()
	assert.Error(t, "close() error", closeErr)
	if !strings.Contains(closeErr.Error(), "close migration source") {
		t.Fatalf("unexpected error: %v", closeErr)
	}
	assert.Equal(t, "sourceDriver.closeCount", sourceDriver.closeCount, 1)
	assert.Equal(t, "databaseDriver.closeCount", databaseDriver.closeCount, 0)
}

func TestSpecNewInstanceCloseReturnsBothErrorsWhenKeepDatabaseOpenFalse(t *testing.T) {
	t.Parallel()

	sourceDriver := &fakeSourceDriver{
		closeErr: errors.New("source close failed"),
	}
	databaseDriver := &fakeDatabaseDriver{
		closeErr: errors.New("database close failed"),
	}
	spec := Spec{
		SourceName: "fake-source",
		NewSourceDriver: func() (source.Driver, error) {
			return sourceDriver, nil
		},
		DatabaseName: "fake-db",
		NewDatabaseDriver: func() (database.Driver, error) {
			return databaseDriver, nil
		},
		KeepDatabaseOpen: false,
	}

	m, close, err := spec.NewInstance()
	assert.NilError(t, "NewInstance() error", err)
	assert.NotNil(t, "NewInstance()", m)
	if close == nil {
		t.Fatalf("close() = nil, want non-nil")
	}

	closeErr := close()
	assert.Error(t, "close() error", closeErr)
	if !strings.Contains(closeErr.Error(), "close migration source") {
		t.Fatalf("unexpected error: %v", closeErr)
	}
	if !strings.Contains(closeErr.Error(), "close migration database") {
		t.Fatalf("unexpected error: %v", closeErr)
	}
	assert.Equal(t, "sourceDriver.closeCount", sourceDriver.closeCount, 1)
	assert.Equal(t, "databaseDriver.closeCount", databaseDriver.closeCount, 1)
}
