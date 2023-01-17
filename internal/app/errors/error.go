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

func WrapAddrNotFoundErr(err error) Error {
	return vsErr{
		code: ErrAddrNotFound,
		msg:  err.Error(),
	}
}

func WrapBadRequest(err error) Error {
	return vsErr{
		code: ErrInvalidParams,
		msg:  err.Error(),
	}
}

func WrapLcdNodeErr(errMsg string) Error {
	return vsErr{
		code: ErrLcdNodeError,
		msg:  errMsg,
	}
}

func WrapNoDataErr() Error {
	return vsErr{
		code: ErrNoData,
		msg:  "no data",
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
