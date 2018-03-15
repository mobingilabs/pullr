package gova

import (
	"bytes"
	"fmt"
)

// ValidationErrors is an error type as collection of validation errors
type ValidationErrors []ValidationError

// Error reports combined validation errors
func (errs ValidationErrors) Error() string {
	buff := bytes.NewBuffer(nil)
	for i := range errs {
		buff.WriteString(errs[i].Field)
		buff.WriteString(": ")
		buff.WriteString(errs[i].Message)
		buff.WriteString("; ")
	}

	return buff.String()
}

// ValidationError is data validation error type
type ValidationError struct {
	Field   string
	Message string
}

// Error return the error message
func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}
