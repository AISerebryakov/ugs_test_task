package http

import (
	"github.com/francoispqt/gojay"
	"github.com/pretcat/ugc_test_task/logger"
	"net/http"
)

const (
	ErrorKey   = "error"
	WarningKey = "warning"
	DataKey    = "data"
)

var nullData = gojay.EmbeddedJSON(`null`)

type Response struct {
	rw         http.ResponseWriter
	statusCode int
	reqId      string
	data gojay.EmbeddedJSON
	err  Error
	warn Warning
}

func NewResponse(rw http.ResponseWriter, reqId string) *Response {
	return &Response{
		rw:         rw,
		reqId:      reqId,
		statusCode: http.StatusOK,
		data:       nullData,
	}
}

func (res *Response) SetData(data []byte) {
	if !res.err.IsEmpty() {
		res.data = nullData
		return
	}
	res.data = data
}

func (res *Response) SetError(err Error) {
	res.err = err
	res.statusCode = err.httpCode
	res.warn.MakeEmpty()
	res.data = nullData
}

func (res *Response) SetWarning(warn Warning) {
	if !res.err.IsEmpty() {
		res.data = nullData
		return
	}
	res.warn = warn
}

func (res Response) Error() Error {
	return res.err
}

func (res Response) Warning() Warning {
	return res.warn
}

func (res *Response) WriteBody() {
	enc := gojay.BorrowEncoder(res.rw)
	defer enc.Release()
	err := enc.EncodeObject(res)

	if err != nil {
		logger.TraceId(res.reqId).AddMsg("error while marshaling body to json").Error(err.Error())
		res.rw.WriteHeader(http.StatusInternalServerError)
		res.rw.Write(encodeResponseErrorJson)
		return
	}
}

func (res *Response) writeHeaders() {
	res.rw.Header().Set(ContentTypeKey, ApplicationJsonKey)
	res.rw.Header().Set(RequestIdKey, res.reqId)
	res.rw.WriteHeader(res.statusCode)
}

func (res *Response) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddObjectKeyNullEmpty(ErrorKey, res.err)
	enc.AddObjectKeyNullEmpty(WarningKey, res.warn)
	enc.AddEmbeddedJSONKey(DataKey, &res.data)
}

func (res *Response) IsNil() bool {
	return res == nil
}
