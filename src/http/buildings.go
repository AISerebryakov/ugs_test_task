package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"ugc_test_task/src/errors"
	"ugc_test_task/src/logger"
	"ugc_test_task/src/managers"
	"ugc_test_task/src/managers/buildings"
	models2 "ugc_test_task/src/models"
)

const (
	buildingsPath = "/v1/buildings"
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
	buildings := make([]models2.Building, 0)
	objectCounter := 0
	err := api.buildingMng.GetBuildings(query, func(building models2.Building) error {
		objectCounter++
		buildings = append(buildings, building)
		return nil
	})
	if err != nil {
		//todo: handle error
		//todo: add details to error
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
		//todo: handle error
		//todo: add details to error
		res.SetError(NewApiError(err))
		return
	}
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

func newGetBuildingsQuery(req Request) (query buildings.GetQuery) {
	urlQuery := req.URL.Query()
	query.ReqId = req.Id()
	query.Id = urlQuery.Get(models2.IdKey)
	query.Address = urlQuery.Get(models2.AddressKey)
	query.FromDate, _ = strconv.ParseInt(urlQuery.Get(managers.FromDateKey), 10, 0)
	query.ToDate, _ = strconv.ParseInt(urlQuery.Get(managers.ToDateKey), 10, 0)
	query.Limit = parseLimit(urlQuery)
	return query
}

func newAddBuildingQuery(req Request) (buildings.AddQuery, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return buildings.AddQuery{}, errors.BodyReadErr.New(err.Error())
	}
	if len(body) == 0 {
		return buildings.AddQuery{}, errors.BodyIsEmpty.New("")
	}
	query, err := buildings.NewAddQueryFromJson(body)
	if err != nil {
		return buildings.AddQuery{}, errors.QueryParseErr.New(err.Error())
	}
	query.ReqId = req.Id()
	return query, nil
}
