package http

import (
	"net/http"
	"net/url"
	"strconv"
	"ugc_test_task/src/common/random"
)

type Request struct {
	*http.Request
	id string
}

func NewRequest(httpReq *http.Request) (req Request) {
	req.Request = httpReq
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

func parseLimit(query url.Values) int {
	limit, _ := strconv.Atoi(query.Get(LimitKey))
	if limit > maxGettingObjects || limit == 0 {
		limit = maxGettingObjects
	}
	return limit
}
