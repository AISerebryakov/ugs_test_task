package http

import (
	"net/http"
	"time"

	"github.com/pretcat/ugc_test_task/common/random"
)

type Request struct {
	*http.Request
	id        string
	startTime time.Time
}

func NewRequest(httpReq *http.Request) (req Request) {
	req.Request = httpReq
	req.startTime = time.Now()
	return req
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

func (req Request) Time() time.Time {
	return req.startTime
}
