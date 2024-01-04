package helpers

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInternal          = errors.New("Something went wrong. Please try again later.")
	ErrInvalidCreditCard = errors.New("Credit card data invalid.")
	ErrUserNotFound      = errors.New("User not found.")
	ErrEmailUsed         = errors.New("User with provided email already exists.")
)

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

func (e ResponseError) Unwrap() error {
	return e.msg
}

func (e ResponseError) Error() string {
	return e.msg.Error()
}

func (e ResponseError) Code() int {
	return e.code
}

type ValidationError struct {
	msg error
}

func NewValidationError(msg error) ValidationError {
	return ValidationError{
		msg: msg,
	}
}

func (e ValidationError) Error() string {
	return e.msg.Error()
}

func (e ValidationError) ErrSlice() []string {
	errors := make([]string, 0)
	for _, err := range strings.Split(e.Error(), ";") {
		errors = append(errors, fmt.Sprintf("Please provide %s field.", strings.Split(err, ":")[0]))
	}

	return errors
}
