package apierror

type FieldError struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func NewValidationError(fieldErrors []FieldError) ValidationError {
	return ValidationError{Errors: fieldErrors}
}
