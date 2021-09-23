package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"ugc_test_task/companymng"
	"ugc_test_task/logger"
	"ugc_test_task/managers"
	"ugc_test_task/models"
)

const (
	companiesPath = "/v1/companies"
)

func (api Api) companyHandlers(res *Response, req Request) {
	switch req.Method {
	//case http.MethodPost:
	//	api.addCompany(rw, req)
	case http.MethodGet:
		api.getCompanies(res, req)
	}
}

func (api Api) getCompanies(res *Response, req Request) {
	query := getCompaniesQuery(req)
	companies := make([]models.Company, 0)
	objectCounter := 0
	err := api.companyMng.GetCompanies(query, func(company models.Company) error {
		objectCounter++
		companies = append(companies, company)
		return nil
	})
	if err != nil {
		//todo: handle error
		//todo: add details to error
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(companies)
	if err != nil {
		apiErr := NewEncodingJsonError("error on encoding companies to json")
		logger.TraceId(req.Id()).AddMsg(apiErr.msg).Error(err.Error())
		res.SetError(apiErr)
		return
	}
	res.SetData(jsonData)
	if objectCounter >= maxGettingObjects {
		res.SetWarning(NewLimitExceededWarning())
	}
}

func (api Api) addCompany(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Add firms"))
}

func getCompaniesQuery(req Request) (query companymng.GetQuery) {
	urlQuery := req.URL.Query()
	query.Id = urlQuery.Get(models.IdKey)
	query.BuildingId = urlQuery.Get(models.BuildingIdKey)
	query.Categories = urlQuery.Get(models.CategoriesKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit, _ = strconv.Atoi(urlQuery.Get(managers.ToDateKey))
	return query
}
