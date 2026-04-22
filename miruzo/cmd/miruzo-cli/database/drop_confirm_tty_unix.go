//go:build !windows

package database

import "os"

func openDropConfirmationIO() (dropConfirmationIO, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return dropConfirmationIO{}, err
	}

	return dropConfirmationIO{
		Reader:  tty,
		Writer:  tty,
		closeFn: tty.Close,
	}, nil
}
