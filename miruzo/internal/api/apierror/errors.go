package apierror

import "sort"

const (
	fieldErrorTypeUnsupported string = "unsupported"
	FieldErrorTypeDuplicate   string = "duplicate"
	FieldErrorTypeInvalid     string = "invalid"
)

type FieldError struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type FieldErrors []FieldError

func (e FieldErrors) Len() int {
	return len(e)
}

func (e FieldErrors) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e FieldErrors) Less(i, j int) bool {
	return e[i].Path < e[j].Path
}

func (e FieldErrors) Sort() {
	sort.Sort(e)
}

func NewUnsupportedError(key string) FieldError {
	return FieldError{
		Type:    fieldErrorTypeUnsupported,
		Path:    "query." + key,
		Message: "is not supported",
	}
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func NewValidationError(fieldErrors []FieldError) ValidationError {
	return ValidationError{Errors: fieldErrors}
}
