package http

import (
	"fmt"
	"net/http"
	"ugc_test_task/src/logger"
)

const (
	RequestIdKey       = "X-Request-Id"
	ApplicationJsonKey = "application/json"
	ContentTypeKey     = "Content-Type"
	LimitKey           = "limit"

	maxGettingObjects = 200
)

type handler struct {
	Api
}

type Api struct {
	server      *http.Server
	conf        Config
	companyMng  CompanyManager
	buildingMng BuildingManager
	categoryMng CategoryManager
}

func NewApi(conf Config) (api Api) {
	api.conf = conf
	api.companyMng = conf.CompanyManager
	api.buildingMng = conf.BuildingManager
	api.categoryMng = conf.CategoryManager
	return api
}

func (api Api) Start(f func(error)) {
	if err := api.startServer(); err != nil {
		f(fmt.Errorf("start http server: %v", err))
	}
}

//todo: add metrics server
//todo: add debug server
//todo: add max objects for getting

func (api *Api) startServer() error {
	if err := api.conf.Validate(); err != nil {
		return fmt.Errorf("config is invalid: %v", err)
	}
	conf := api.conf
	api.server = &http.Server{
		Addr:              conf.Address(),
		Handler:           handler{*api},
		ReadTimeout:       conf.ReadTimeout.TimeDuration(),
		ReadHeaderTimeout: conf.ReadHeaderTimeout.TimeDuration(),
		WriteTimeout:      conf.WriteTimeout.TimeDuration(),
		IdleTimeout:       conf.IdleTimeout.TimeDuration(),
		MaxHeaderBytes:    conf.MaxHeaderBytes.Int(),
	}
	logger.Msg("start http server").Info(conf.Address())
	if err := api.server.ListenAndServe(); err != nil {
		//todo: handle error
		return err
	}
	return nil
}

func (h handler) ServeHTTP(rw http.ResponseWriter, httpReq *http.Request) {
	req := NewRequest(httpReq)
	res := NewResponse(rw, req.Id())
	fmt.Println(req.Path())

	switch req.Path() {
	case companiesPath:
		h.companyHandlers(res, req)
	case buildingsPath:
		h.buildingHandlers(res, req)
	case categoriesPath:
		h.categoriesHandlers(res, req)
	}
	if !res.err.IsEmpty() {
		logger.TraceId(req.Id()).Error(res.Error().Error())
	}
	res.writeHeaders()
	res.WriteBody()
}
