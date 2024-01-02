package helpers

import "errors"

var ErrInternal = errors.New("Something went wrong. Please try again later.")

type ResponseError struct {
	msg  error
	code int
}

func NewResponseError(msg error, code int) ResponseError {
	return ResponseError{
		msg:  msg,
		code: code,
	}
}

func (e ResponseError) Error() string {
	return e.msg.Error()
}

func (e ResponseError) Code() int {
	return e.code
}
