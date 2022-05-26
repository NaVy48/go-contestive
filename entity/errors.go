package entity

import "fmt"

type errorWrapper struct {
	message  string
	innerErr error
}

func (err errorWrapper) Error() string {
	if err.innerErr != nil {
		return fmt.Sprintf("%s: %s", err.message, err.innerErr.Error())
	} else {
		return err.message
	}
}
func (err errorWrapper) Unwrap() error {
	return err.innerErr
}

func (err errorWrapper) Is(target error) bool {
	t, ok := target.(errorWrapper)
	return ok && t.message == err.message
}
func (err errorWrapper) As(target interface{}) bool {
	t, ok := target.(*errorWrapper)
	if ok && t.message == err.message {
		*t = err
		return true
	}
	return false
}

func ErrCustomWrapper(message string, innerErr error) error {
	return errorWrapper{message, innerErr}
}

func ErrNotFound(innerErr error) error {
	return errorWrapper{"entity not found", innerErr}
}

func ErrBadCredentials(innerErr error) error {
	return errorWrapper{"incorect username or password", innerErr}
}

func ErrInvalidToken(innerErr error) error {
	return errorWrapper{"invalid token", innerErr}
}

func ErrDataConflict(innerErr error) error {
	return errorWrapper{"data conflict", innerErr}
}

func ErrPackageNotCompatible(innerErr error) error {
	return errorWrapper{"package is not compatible with current problem", innerErr}
}

func ErrRevisionTooOld(innerErr error) error {
	return errorWrapper{"problem has newer revision", innerErr}
}

type ErrValidation struct {
	fieldName string
	rule      string
	object    interface{}
}

// NewValidationError creates ValidationError
// example: NewValidationError("Name", "is required", struct{Name string}{""})
func NewValidationError(fieldName, rule string, object interface{}) ErrValidation {
	return ErrValidation{fieldName, rule, object}
}

func (err ErrValidation) Error() string {
	return fmt.Sprintf("%s %s in %v", err.fieldName, err.rule, err.object)
}
