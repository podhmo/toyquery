package core

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorFactory :
type ErrorFactory struct {
	reason error
}

// Reason :
func (f *ErrorFactory) Reason() error {
	return f.reason
}

// Error :
func (f *ErrorFactory) Error() string {
	return f.reason.Error()
}

// New :
func (f *ErrorFactory) New(msg string) error {
	return errors.WithMessage(
		&Error{reason: f.reason, raw: f.reason},
		msg,
	)
}

// Wrap :
func (f *ErrorFactory) Wrap(err error, msg string) error {
	return errors.WithMessage(
		&Error{reason: f.reason, raw: err},
		msg,
	)
}

// Error :
type Error struct {
	raw    error
	reason error
}

// Error :
func (e *Error) Error() string {
	if e.raw == e.reason {
		return e.raw.Error()
	}
	return fmt.Sprintf("%s (reason=%s)", e.raw.Error(), e.reason.Error())
}

// Cause :
func (e *Error) Cause() error {
	return e.raw
}

// Reason :
func (e *Error) Reason() error {
	return e.reason
}

// DefineError :
func DefineError(reason error) *ErrorFactory {
	return &ErrorFactory{reason: reason}
}

// TODO: more gentle implementation
var (
	// ErrDatabaseNotFound :
	ErrDatabaseNotFound = DefineError(fmt.Errorf("database not found"))
	// ErrTableNotFound :
	ErrTableNotFound = DefineError(fmt.Errorf("table not found"))
	// ErrFieldNotFound :
	ErrFieldNotFound = DefineError(fmt.Errorf("field not found"))

	// ErrRecordNotFound :
	ErrRecordNotFound = DefineError(fmt.Errorf("record not found"))

	// ErrSomethingWrong :
	ErrSomethingWrong = DefineError(fmt.Errorf("something wrong")) // ???
)
