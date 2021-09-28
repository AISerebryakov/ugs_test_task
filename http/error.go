package http

import (
	"net/http"
	"strconv"

	"github.com/francoispqt/gojay"
	"github.com/pretcat/ugc_test_task/errors"
)

const (
	InternalServerErrorTitle = "internal_server_error"
	IncorrectRequestTitle    = "incorrect_request"
	EncodingJsonErrorTitle   = "encoding_json_error"

	TitleKey = "title"
	MsgKey   = "msg"
)

var (
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
	httpCode int
	title    string
	msg      string
}

func NewApiError(err error) Error {
	errType := errors.GetType(err)
	switch errType {
	case errors.QueryIsInvalid, errors.QueryParseErr,
		errors.BodyReadErr, errors.BodyIsEmpty,
		errors.Duplicate, errors.InputParamsIsInvalid:
		return NewIncorrectRequestError(err.Error())
	default:
		return NewInternalServerError(err.Error())
	}
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

func (err Error) Error() string {
	return err.msg
}

func (err Error) String() string {
	return strconv.Itoa(err.httpCode) + " " + err.title + " " + err.msg
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
