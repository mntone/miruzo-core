package database

import "io"

type dropConfirmationIO struct {
	Reader io.Reader
	Writer io.Writer

	closeFn func() error
}

func (c dropConfirmationIO) Close() error {
	if c.closeFn == nil {
		return nil
	}
	return c.closeFn()
}

var openDropConfirmationIOFn = openDropConfirmationIO
