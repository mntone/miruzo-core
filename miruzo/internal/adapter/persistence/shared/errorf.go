package shared

import (
	"fmt"
	"strings"
)

func WrapKV(base error, operation string, keyValues ...any) error {
	if base == nil {
		panic("WrapKV: base must not be nil")
	}
	if operation == "" {
		panic("WrapKV: operation must not be empty")
	}
	if len(keyValues)%2 != 0 {
		panic("WrapKV: keyValues must be key/value pairs")
	}

	var b strings.Builder
	b.WriteString(": operation=")
	b.WriteString(operation)

	for i := 0; i < len(keyValues); i += 2 {
		key, ok := keyValues[i].(string)
		if !ok || key == "" {
			panic("WrapKV: keys must be non-empty string")
		}

		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteByte('=')
		fmt.Fprint(&b, keyValues[i+1])
	}

	return fmt.Errorf("%w%s", base, b.String())
}
