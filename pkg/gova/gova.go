package gova

import (
	"fmt"
	"strings"
)

// Validator is an helper for validating struct fields
type Validator struct {
	errors []ValidationError
}

// Valid reports true if no errors encountered
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// Errors reports validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// Extend extends the validation errors with another validator's errors.
// New errors will have field names prefixed with given field.
func (v *Validator) Extend(field string, other ValidationErrors) {
	for _, err := range other {
		field := fmt.Sprintf("%s.%s", field, err.Field)
		v.errors = append(v.errors, ValidationError{field, err.Message})
	}
}

// ExtendElt extends the validation errors with array element's validation errors.
// It is useful if you want to use extend with each element in your array.
func (v *Validator) ExtendElt(field string, index int, other ValidationErrors) {
	field = fmt.Sprintf("%s[%d]", field, index)
	v.Extend(field, other)
}

// ShouldBeOneOf checks if the given string is equal to one of the allowed values
func (v *Validator) ShouldBeOneOf(field string, value string, allowed ...string) {
	for _, str := range allowed {
		if value == str {
			return
		}
	}

	msg := fmt.Sprintf(`should be one of the "%s"`, strings.Join(allowed, ", "))
	v.errors = append(v.errors, ValidationError{field, msg})
}

// NotEmptyString checks if the given string is not empty
func (v *Validator) NotEmptyString(field string, value string) {
	if value == "" {
		v.errors = append(v.errors, ValidationError{field, "can not be empty"})
	}
}

// NotEmpty checks if the given slice is not empty
func (v *Validator) NotEmpty(field string, nitems int) {
	if nitems == 0 {
		v.errors = append(v.errors, ValidationError{field, "should have at least one value"})
	}
}

// NonZero checks if the given integer is not zero
func (v *Validator) NonZero(field string, value int) {
	if value == 0 {
		v.errors = append(v.errors, ValidationError{field, "can not be zero"})
	}
}
