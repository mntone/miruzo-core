package bind

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func bitSizeOfSignedInteger[T constraints.Signed]() int {
	var zero T
	switch any(zero).(type) {
	case int8:
		return 8
	case int16:
		return 16
	case int32:
		return 32
	case int64:
		return 64
	case int:
		return strconv.IntSize
	default:
		return 64
	}
}
