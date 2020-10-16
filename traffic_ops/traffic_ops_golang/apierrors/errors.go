package apierrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Errors represents a set of errors to be handled by the API.
type Errors struct {
	// Code is the HTTP response code. If no error has occurred, this *should*
	// be http.StatusOK.
	Code int
	// SystemError is an error that occurred internally; not safe for exposure
	// to the user.
	SystemError error
	// UserError is an error that should be shown to the user to explain why
	// their request failed.
	UserError error
}

// New constructs a new Errors where no error has actually occurred.
func New() Errors {
	return Errors{
		Code:        http.StatusOK,
		SystemError: nil,
		UserError:   nil,
	}
}

// Occurred returns whether at least one error has occurred (is non-nil).
func (e Errors) Occurred() bool {
	return e.SystemError != nil || e.UserError != nil
}

// SetSystemError sets the Errors' system-level error to a new error containing
// the passed message.
func (e *Errors) SetSystemError(err string) {
	e.SystemError = errors.New(err)
}

// SetUserError sets the Errors' user-level error to a new error containing the
// passed message.
func (e *Errors) SetUserError(err string) {
	e.UserError = errors.New(err)
}

func (e Errors) String() string {
	return fmt.Sprintf("Errors(Code=%d, SystemError='%v', UserError='%v')", e.Code, e.SystemError, e.UserError)
}
