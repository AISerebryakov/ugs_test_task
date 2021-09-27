package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"ugc_test_task/errors"
	"ugc_test_task/logger"
	"ugc_test_task/managers"
	categmng "ugc_test_task/managers/categories"
	"ugc_test_task/models"
)

const (
	SearchByNameKey = "search_by_name"

	categoriesPath = "/v1/categories"
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
	fmt.Println("Query: ", query)
	categories := make([]models.Category, 0)
	objectCounter := 0
	err := api.categoryMng.GetCategories(query, func(category models.Category) error {
		objectCounter++
		categories = append(categories, category)
		return nil
	})
	if err != nil {
		//todo: handle error
		//todo: add details to error
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
		//todo: handle error
		//todo: add details to error
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

func newGetCategoriesQuery(req Request) (query categmng.GetQuery) {
	urlQuery := req.URL.Query()
	query.ReqId = req.Id()
	query.Id = urlQuery.Get(models.IdKey)
	query.Name = urlQuery.Get(models.NameKey)
	query.SearchName = urlQuery.Get(SearchByNameKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	return query
}

func newAddCategoryQuery(req Request) (categmng.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return categmng.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return categmng.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := categmng.NewAddQueryFromJson(body)
	if err != nil {
		return categmng.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	query.ReqId = req.Id()
	return query, nil
}
