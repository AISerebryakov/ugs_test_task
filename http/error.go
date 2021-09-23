package http

import (
	"errors"
	"net/http"
	"ugc_test_task/managers"

	"github.com/francoispqt/gojay"
)

const (
	InternalServerErrorTitle = "internal_server_error"
	IncorrectRequestTitle    = "incorrect_request"
	EncodingJsonErrorTitle   = "encoding_json_error"

	TitleKey = "title"
	MsgKey   = "msg"
)

var (
	ErrBodyIsEmpty = errors.New("body is empty")
	ErrBodyReading = errors.New("body reading")

	encodeResponseErrorJson = []byte(`{
	"data": null,
	"error": {
		"title": "encoding_json_error",
		"msg": "error on encoding response to json"
	},
	"warning": null
}`)
)

type Error struct {
	err      error
	httpCode int
	title    string
	msg      string
}

func NewApiError(err error) Error {
	if apiErr, ok := matchManagerErrors(err); ok {
		apiErr.err = err
		return apiErr
	}
	if apiErr, ok := matchApiErrors(err); ok {
		apiErr.err = err
		return apiErr
	}
	apiErr := NewInternalServerError(err.Error())
	apiErr.err = err
	return apiErr
}

func matchManagerErrors(err error) (Error, bool) {
	if errors.Is(err, managers.ErrQueryInvalid) {
		return NewIncorrectRequestError(err.Error()), true
	}
	if errors.Is(err, managers.ErrParsingQuery) {
		return NewIncorrectRequestError(err.Error()), true
	}
	if errors.Is(err, managers.ErrSaveToDb) {
		return NewInternalServerError(managers.ErrSaveToDb.Error()), true
	}
	return Error{}, false
}

func matchApiErrors(err error) (Error, bool) {
	if errors.Is(err, ErrBodyIsEmpty) {
		return NewIncorrectRequestError(ErrBodyIsEmpty.Error()), true
	}
	if errors.Is(err, ErrBodyReading) {
		return NewIncorrectRequestError(err.Error()), true
	}
	return Error{}, false
}

func NewIncorrectRequestError(msg string) Error {
	return Error{
		httpCode: http.StatusBadRequest,
		title:    IncorrectRequestTitle,
		msg:      msg,
	}
}

func NewInternalServerError(msg string) Error {
	return Error{
		httpCode: http.StatusInternalServerError,
		title:    InternalServerErrorTitle,
		msg:      msg,
	}
}

func NewEncodingJsonError(msg string) Error {
	return Error{
		httpCode: http.StatusInternalServerError,
		title:    EncodingJsonErrorTitle,
		msg:      msg,
	}
}

func (err Error) OriginError() error {
	return err.err
}

func (err Error) Error() string {
	return err.msg
}

func (err Error) IsEmpty() bool {
	return len(err.title) == 0
}

func (err Error) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey(TitleKey, err.title)
	enc.AddStringKey(MsgKey, err.msg)
}

func (err Error) IsNil() bool {
	return err.IsEmpty()
}
