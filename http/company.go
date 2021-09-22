package http

import (
	"fmt"
	"net/http"
	"ugc_test_task/companymng"
	"ugc_test_task/models"
)

const (
	companiesPath = "/v1/companies"
)

func (api Api) companyHandlers(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		api.addCompany(rw, req)
	case http.MethodGet:
		api.getCompanies(rw, req)
	}
}

func (api Api) getCompanies(rw http.ResponseWriter, req *http.Request) {
	ids := req.URL.Query()[models.IdKey]
	id := ""
	if len(ids) > 0 {
		id = ids[0]
	}
	fmt.Println("Id: ", id)
	query := companymng.GetQuery{Id: id}
	err := api.companyMng.GetCompanies(query, func(company models.Company) error {
		fmt.Println("Company: ", company)
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(company.Name))
		return nil
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		//todo: handle error
	}
	//todo: check limit
}

func (api Api) addCompany(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Add firms"))
}
