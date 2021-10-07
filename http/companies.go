package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/logger"
	"github.com/pretcat/ugc_test_task/managers"
	compmng "github.com/pretcat/ugc_test_task/managers/companies"
	"github.com/pretcat/ugc_test_task/models"
)

const (
	SearchByCategoryKey = "search_by_category"
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
	logger.TraceId(req.Id()).AddMsg("query").Debugf("%#v", query)
	companies := make([]models.Company, 0)
	objectCounter := 0
	err := api.companyMng.GetCompanies(query, func(company models.Company) error {
		objectCounter++
		companies = append(companies, company)
		return nil
	})
	if err != nil {
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
		res.SetError(NewApiError(err))
		return
	}
	logger.TraceId(req.Id()).AddMsg("query").Debugf("%#v", query)
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

func newGetCompaniesQuery(req Request) (query compmng.GetQuery) {
	urlQuery := req.URL.Query()
	query.TraceId = req.Id()
	query.Id = urlQuery.Get(models.IdKey)
	query.Category = urlQuery.Get(SearchByCategoryKey)
	query.BuildingId = urlQuery.Get(models.BuildingIdKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	query.Offset = parseOffset(urlQuery)
	query.Ascending.Exists, query.Ascending.Value = parseAscending(urlQuery)
	return query
}

func newAddCompanyQuery(req Request) (compmng.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return compmng.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return compmng.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := compmng.NewAddQueryFromJson(body)
	if err != nil {
		return compmng.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	query.TraceId = req.Id()
	return query, nil
}
