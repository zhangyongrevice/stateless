package serr

import (
	"github.com/pkg/errors"
)

type cerr struct {
	error
	code uint32
}

// Cause make cerr plays nice with errors.Cause() by returning the underlying
// error
func (e *cerr) Cause() error {
	return e.error
}

// New creates new error containing code and msg
func New(msg string, code uint32) error {
	return &cerr{
		error: errors.New(msg),
		code:  code,
	}
}

// Wrap creates new error containing the cause, code, and msg
// if err is nil, will return nil
func Wrap(err error, msg string, code uint32) error {
	if err == nil {
		return nil
	}

	if msg == "" {
		return &cerr{
			error: err,
			code:  code,
		}
	}

	return &cerr{
		error: errors.Wrap(err, msg),
		code:  code,
	}
}

// Code returns the error code of the given error.
// It's common for low level error to be wrapped within higher level message.
// So, this function will recursively go deeper to find and return the code of
// first cerr in the chain.
//
// Note that this only work for error created by either cerr.Wrap() or errors.Wrap()
func Code(err error) uint32 {
	if err == nil {
		return uint32(Constant_SUCCESS)
	}

	if cerr, ok := err.(*cerr); ok {
		return cerr.code
	}

	type causer interface {
		Cause() error
	}

	cause, ok := err.(causer)
	if !ok {
		return uint32(Constant_ERROR_UNKNOW)
	}
	return Code(cause.Cause())
}
