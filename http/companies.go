package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"ugc_test_task/errors"
	"ugc_test_task/logger"
	"ugc_test_task/managers"
	companmng "ugc_test_task/managers/companies"
	"ugc_test_task/models"
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

func newGetCompaniesQuery(req Request) (query companmng.GetQuery) {
	urlQuery := req.URL.Query()
	query.ReqId = req.Id()
	query.Id = urlQuery.Get(models.IdKey)
	query.BuildingId = urlQuery.Get(models.BuildingIdKey)
	query.Categories = urlQuery.Get(models.CategoriesKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	return query
}

func newAddCompanyQuery(req Request) (companmng.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return companmng.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return companmng.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := companmng.NewAddQueryFromJson(body)
	if err != nil {
		return companmng.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	return query, nil
}
