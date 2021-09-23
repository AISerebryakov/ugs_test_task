package http

import (
	"fmt"
	"net/http"
	"ugc_test_task/logger"
)

const (
	RequestIdKey       = "X-Request-Id"
	ApplicationJsonKey = "application/json"
	ContentTypeKey     = "Content-Type"

	maxGettingObjects = 200
)

type handler struct {
	Api
}

type Api struct {
	server     *http.Server
	conf       Config
	companyMng CompanyManager
}

func NewApi(conf Config) (api Api, _ error) {
	if err := conf.Validate(); err != nil {
		return Api{}, fmt.Errorf("config is invalid: %v", err)
	}
	api.conf = conf
	api.companyMng = conf.CompanyManager
	return api, nil
}

func (api Api) Start(f func(error)) {
	if err := api.startServer(); err != nil {
		f(fmt.Errorf("start http server: %v", err))
	}
}

//todo: add metrics server
//todo: add debug server

func (api *Api) startServer() error {
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
		h.companyHandlers(&res, req)
	}
	if !res.err.IsEmpty() {
		logger.Error(res.err.msg)
	}
	fmt.Println("JSON API: ", string(res.data))
	res.writeHeaders()
	res.WriteBody()
}
