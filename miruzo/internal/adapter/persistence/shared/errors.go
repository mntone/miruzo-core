package shared

import "fmt"

func JoinErrors(primary error, secondary error) error {
	if secondary == nil {
		return primary
	}

	return fmt.Errorf("%w (cleanup failed: %v)", primary, secondary)
}
