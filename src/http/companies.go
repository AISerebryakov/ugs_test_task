package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"ugc_test_task/src/errors"
	"ugc_test_task/src/logger"
	"ugc_test_task/src/managers"
	"ugc_test_task/src/managers/companies"
	models2 "ugc_test_task/src/models"
)

const (
	companiesPath = "/v1/companies"
)

func (api Api) companyHandlers(res *Response, req Request) {
	switch req.Method {
	case http.MethodPost:
		api.addCompany(res, req)
	case http.MethodGet:
		api.getCompanies(res, req)
	}
}

func (api Api) getCompanies(res *Response, req Request) {
	query := newGetCompaniesQuery(req)
	companies := make([]models2.Company, 0)
	objectCounter := 0
	err := api.companyMng.GetCompanies(query, func(company models2.Company) error {
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

func (api Api) addCompany(res *Response, req Request) {
	query, err := newAddCompanyQuery(req)
	if err != nil {
		//todo: handle error
		//todo: add details to error
		res.SetError(NewApiError(err))
		return
	}
	comp, err := api.companyMng.AddCompany(query)
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(comp)
	if err != nil {
		apiErr := NewEncodingJsonError("error on encoding company to json")
		logger.TraceId(req.Id()).AddMsg(apiErr.msg).Error(err.Error())
		res.SetError(apiErr)
		return
	}
	res.SetData(jsonData)
}

func newGetCompaniesQuery(req Request) (query companies.GetQuery) {
	urlQuery := req.URL.Query()
	query.ReqId = req.Id()
	query.Id = urlQuery.Get(models2.IdKey)
	query.BuildingId = urlQuery.Get(models2.BuildingIdKey)
	query.Categories = urlQuery.Get(models2.CategoriesKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	return query
}

func newAddCompanyQuery(req Request) (companies.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return companies.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return companies.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := companies.NewAddQueryFromJson(body)
	if err != nil {
		return companies.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	return query, nil
}
