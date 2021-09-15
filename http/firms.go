package http

import "net/http"

func (api Api) firmHandlers(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		api.addFirm(rw, req)
	case http.MethodGet:
		api.getFirms(rw, req)
	}
}

func (api Api) getFirms(rw http.ResponseWriter, req *http.Request) {

}

func (api Api) addFirm(rw http.ResponseWriter, req *http.Request) {

}
