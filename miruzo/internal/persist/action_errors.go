package persist

import "fmt"

var ErrActionAlreadyExists = fmt.Errorf("%w: action already exists", ErrConflict)
