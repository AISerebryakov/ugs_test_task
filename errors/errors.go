package errors

import (
	"errors"
)

type Type string

const (
	EmptyType            Type = ""
	Duplicate            Type = "duplicate"
	QueryIsInvalid       Type = "query is invalid"
	QueryParseErr        Type = "query parse error"
	BodyReadErr          Type = "body read error"
	BodyIsEmpty          Type = "body is empty"
	InputParamsIsInvalid Type = "input parameters is invalid"
)

type Error struct {
	t   Type
	err error
	msg string
}

func New(t Type, err error) *Error {
	return &Error{t: t, err: err, msg: err.Error()}
}

func GetType(err error) Type {
	var e *Error
	if errors.As(err, &e) {
		return e.t
	}
	return EmptyType
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Wrap(err error, msg string) error {
	if len(msg) == 0 || err == nil {
		return err
	}
	var e *Error
	if errors.As(err, &e) {
		e.msg = msg + ": " + e.msg
		return e
	}
	return New(EmptyType, err).AddBefore(msg)
}

func (err Error) Type() Type {
	return err.t
}

func (err Error) Error() string {
	return err.msg
}

func (err *Error) Add(msg string) *Error {
	if len(msg) == 0 {
		return err
	}
	err.msg = err.msg + ": " + msg
	return err
}

func (err *Error) AddBefore(msg string) *Error {
	if len(msg) == 0 {
		return err
	}
	err.msg = msg + ": " + err.msg
	return err
}

func (t Type) New(msg string) *Error {
	if len(msg) == 0 {
		return &Error{t: t, msg: string(t)}
	}
	return &Error{t: t, msg: string(t) + ": " + msg}
}
