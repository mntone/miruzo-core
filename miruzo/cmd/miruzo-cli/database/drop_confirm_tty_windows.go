package database

import (
	"errors"
	"os"
)

func openDropConfirmationIO() (dropConfirmationIO, error) {
	conIn, err := os.OpenFile("CONIN$", os.O_RDONLY, 0)
	if err != nil {
		return dropConfirmationIO{}, err
	}

	conOut, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
	if err != nil {
		_ = conIn.Close()
		return dropConfirmationIO{}, err
	}

	return dropConfirmationIO{
		Reader: conIn,
		Writer: conOut,
		closeFn: func() error {
			return errors.Join(conOut.Close(), conIn.Close())
		},
	}, nil
}
