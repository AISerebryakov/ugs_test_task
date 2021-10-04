package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/logger"
	"github.com/pretcat/ugc_test_task/managers"
	categories2 "github.com/pretcat/ugc_test_task/managers/categories"
	"github.com/pretcat/ugc_test_task/models"
)

func (api Api) categoriesHandlers(res *Response, req Request) {
	switch req.Method {
	case http.MethodPost:
		api.addCategory(res, req)
	case http.MethodGet:
		api.getCategories(res, req)
	}

}

func (api Api) getCategories(res *Response, req Request) {
	query := newGetCategoriesQuery(req)
	logger.TraceId(req.Id()).AddMsg("query").Debug(fmt.Sprintf("%#v", query))
	categories := make([]models.Category, 0)
	objectCounter := 0
	err := api.categoryMng.GetCategories(query, func(category models.Category) error {
		objectCounter++
		categories = append(categories, category)
		return nil
	})
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(categories)
	if err != nil {
		res.SetError(NewEncodingJsonError("error on encoding categories to json"))
		return
	}
	res.SetData(jsonData)
	if objectCounter >= maxGettingObjects {
		res.SetWarning(NewLimitExceededWarning())
	}
}

func (api Api) addCategory(res *Response, req Request) {
	query, err := newAddCategoryQuery(req)
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	comp, err := api.categoryMng.AddCategory(query)
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(comp)
	if err != nil {
		apiErr := NewEncodingJsonError("error on encoding category to json")
		logger.TraceId(req.Id()).AddMsg(apiErr.msg).Error(err.Error())
		res.SetError(apiErr)
		return
	}
	res.SetData(jsonData)
}

func newGetCategoriesQuery(req Request) (query categories2.GetQuery) {
	urlQuery := req.URL.Query()
	query.TraceId = req.Id()
	query.Id = urlQuery.Get(models.IdKey)
	query.SetName(urlQuery.Get(SearchByNameKey))
	query.SetNameStrict(urlQuery.Get(SearchByNameStrictKey))
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Offset, _ = strconv.Atoi(urlQuery.Get(OffsetKey))
	query.Ascending.Exists, query.Ascending.Value = parseAscending(urlQuery)
	query.Limit = parseLimit(urlQuery)
	return query
}

func newAddCategoryQuery(req Request) (categories2.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return categories2.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return categories2.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := categories2.NewAddQueryFromJson(body)
	if err != nil {
		return categories2.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	query.ReqId = req.Id()
	return query, nil
}
