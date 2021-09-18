package http

import (
	"net/http"
	"ugc_test_task/firmmng"
	"ugc_test_task/models"
)

const (
	firmPath = "/api/v1/firm"
)

func (api Api) firmHandlers(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		api.addFirm(rw, req)
	case http.MethodGet:
		api.getFirms(rw, req)
	}
}

func (api Api) getFirms(rw http.ResponseWriter, req *http.Request) {
	query := firmmng.GetQuery{Id: "Test id"}
	api.firmMng.GetFirms(query, func(firm models.Company) error {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(firm.Name))
		return nil
	})
}

func (api Api) addFirm(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Add firms"))
}
