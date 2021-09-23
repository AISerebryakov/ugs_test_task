package http

import (
	"net/http"
	"ugc_test_task/common/random"
)

type Request struct {
	*http.Request
	id string
}

func NewRequest(httpReq *http.Request) Request {
	return Request{Request: httpReq}
}

func (req *Request) Id() string {
	if len(req.id) > 0 {
		return req.id
	}
	id := req.Header.Get(RequestIdKey)
	if len(id) > 0 {
		req.id = id
		return id
	}
	req.id = random.GenerateRequestId()
	return req.id
}

func (req Request) Path() string {
	return req.URL.Path
}
