package errors

import "fmt"

type Error interface {
	Error() string
	Code() int
	Msg() string
}

// Wrap system error
func Wrap(err error) Error {
	return vsErr{
		code: ErrSystemError,
		msg:  err.Error(),
	}
}

// Wrapf system error
func Wrapf(format string, a ...interface{}) Error {
	return Wrap(fmt.Errorf(format, a))
}

// WrapDetail with code and msg
func WrapDetail(code int, msg string) Error {
	return vsErr{
		code: code,
		msg:  msg,
	}
}

func WrapBadRequest(msg string) Error {
	return vsErr{
		code: ErrBadRequest,
		msg:  msg,
	}
}

func WrapTxNotFound() Error {
	return vsErr{
		code: ErrTxNotFound,
		msg:  "Tx not found",
	}
}

func WrapTxNotUnique(msg string) Error {
	return vsErr{
		code: ErrTxNotUnique,
		msg:  msg,
	}
}

type vsErr struct {
	code int
	msg  string
}

func (e vsErr) Error() string {
	return fmt.Sprintf("err_code: %d, err_msg: %s", e.code, e.msg)
}

func (e vsErr) Code() int {
	return e.code
}

func (e vsErr) Msg() string {
	return e.msg
}
