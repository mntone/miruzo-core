package contract

type BindVarStyle string

const (
	BindVarStyleDollar   BindVarStyle = "$"
	BindVarStyleQuestion BindVarStyle = "?"
)

type DBErrorMapping string

const (
	DBErrorMappingNone    DBErrorMapping = "none"
	DBErrorMappingDefault DBErrorMapping = "default"
	DBErrorMappingDelete  DBErrorMapping = "delete"
)

type Dialect interface {
	MapError(operation string, err error, mapping DBErrorMapping) error
	BindVarStyle() BindVarStyle

	// Param returns placeholder for index.
	Param(index int32) string

	// ParamRange returns placeholders for [start, end] (inclusive).
	ParamRange(start, end int32) []any
}

func ParamRange(start, end int32, param func(index int32) string) []any {
	if start > end {
		panic("invalid param range: start > end")
	}

	length := end - start + 1
	params := make([]any, length)
	for i := int32(0); i < length; i += 1 {
		params[i] = param(start + i)
	}
	return params
}
