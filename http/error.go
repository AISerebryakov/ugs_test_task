package http

import (
	"net/http"

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
	writeResponseErrorJson = []byte(`{
	"data": null,
	"error": {
		"title": "internal_server_error",
		"msg": "error on writing result to response"
	},
	"warning": null
}`)
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
	if apiErr, ok := matchContextErrors(err); ok {
		return apiErr
	}
	return NewInternalServerError(err.Error())
}

func matchContextErrors(err error) (Error, bool) {
	//if errors.Is(err, context.CtxIsInvalidErr) {
	//	return newIncorrectRequestError(err.Error()), true
	//}
	return Error{}, false
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
