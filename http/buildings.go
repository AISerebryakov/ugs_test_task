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
	buildmng "github.com/pretcat/ugc_test_task/managers/buildings"
	"github.com/pretcat/ugc_test_task/models"
)

func (api Api) buildingHandlers(res *Response, req Request) {
	switch req.Method {
	case http.MethodPost:
		api.addBuilding(res, req)
	case http.MethodGet:
		api.getBuildings(res, req)
	}
}

func (api Api) getBuildings(res *Response, req Request) {
	query := newGetBuildingsQuery(req)
	logger.TraceId(req.Id()).AddMsg("query").Debug(fmt.Sprintf("%#v", query))
	buildings := make([]models.Building, 0)
	objectCounter := 0
	err := api.buildingMng.GetBuildings(query, func(building models.Building) error {
		objectCounter++
		buildings = append(buildings, building)
		return nil
	})
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(buildings)
	if err != nil {
		apiErr := NewEncodingJsonError("error on encoding buildings to json")
		logger.TraceId(req.Id()).AddMsg(apiErr.msg).Error(err.Error())
		res.SetError(apiErr)
		return
	}
	res.SetData(jsonData)
	if objectCounter >= maxGettingObjects {
		res.SetWarning(NewLimitExceededWarning())
	}
}

func (api Api) addBuilding(res *Response, req Request) {
	query, err := newAddBuildingQuery(req)
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	logger.TraceId(req.Id()).AddMsg("query").Debugf("%#v", query)
	building, err := api.buildingMng.AddBuilding(query)
	if err != nil {
		res.SetError(NewApiError(err))
		return
	}
	jsonData, err := json.Marshal(building)
	if err != nil {
		apiErr := NewEncodingJsonError("error on encoding building to json")
		logger.TraceId(req.Id()).AddMsg(apiErr.msg).Error(err.Error())
		res.SetError(apiErr)
		return
	}
	res.SetData(jsonData)
}

func newGetBuildingsQuery(req Request) (query buildmng.GetQuery) {
	urlQuery := req.URL.Query()
	query.TraceId = req.Id()
	query.Id = urlQuery.Get(models.IdKey)
	query.Address = urlQuery.Get(models.AddressKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	query.Offset = parseOffset(urlQuery)
	query.Ascending.Exists, query.Ascending.Value = parseAscending(urlQuery)
	return query
}

func newAddBuildingQuery(req Request) (buildmng.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return buildmng.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return buildmng.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := buildmng.NewAddQueryFromJson(body)
	if err != nil {
		return buildmng.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	query.TraceId = req.Id()
	return query, nil
}
