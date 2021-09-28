package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"github.com/pretcat/ugc_test_task/src/errors"
	"github.com/pretcat/ugc_test_task/src/logger"

	"github.com/arl/statsviz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	RequestIdKey       = "X-Request-Id"
	ApplicationJsonKey = "application/json"
	ContentTypeKey     = "Content-Type"
	LimitKey           = "limit"

	maxGettingObjects = 200
)

var (
	newRequestLogStr = func(method, path string) string {
		return fmt.Sprintf("new http request: %s %s", method, path)
	}
	finalRequestLogStr = func(method, path string) string {
		return fmt.Sprintf("final http request: %s %s", method, path)
	}
)

type Api struct {
	server        *http.Server
	debugServer   *http.Server
	metricsServer *http.Server
	conf          Config
	companyMng    CompanyManager
	buildingMng   BuildingManager
	categoryMng   CategoryManager
}

func NewApi(conf Config) *Api {
	api := new(Api)
	api.conf = conf
	api.companyMng = conf.CompanyManager
	api.buildingMng = conf.BuildingManager
	api.categoryMng = conf.CategoryManager
	return api
}

func (api *Api) Start(f func(error)) {
	go func() {
		if err := api.startMetricsServer(); err != nil {
			f(fmt.Errorf("start metrics http server: %v", err))
		}
	}()
	go func() {
		if err := api.startDebugServer(); err != nil {
			f(fmt.Errorf("start debug http server: %v", err))
		}
	}()
	go func() {
		if err := api.startServer(); err != nil {
			f(fmt.Errorf("start http server: %v", err))
		}
	}()
}

func (api *Api) Shutdown(ctx context.Context) {
	if api == nil {
		return
	}
	if api.server != nil {
		if err := api.server.Shutdown(ctx); err != nil {
			logger.Errorf("api http server shutdown: %v", err)
		}
		logger.Info("api http server shutdown")
	}
	if api.metricsServer != nil {
		if err := api.metricsServer.Shutdown(ctx); err != nil {
			logger.Errorf("metrics http server shutdown: %v", err)
		}
		logger.Info("metrics http server shutdown")
	}
	if api.debugServer != nil {
		if err := api.debugServer.Shutdown(ctx); err != nil {
			logger.Errorf("debug http server shutdown: %v", err)
		}
		logger.Info("debug http server shutdown")
	}
}

func (api *Api) startServer() error {
	if err := api.conf.Validate(); err != nil {
		return fmt.Errorf("config is invalid: %v", err)
	}
	conf := api.conf
	api.server = &http.Server{
		Addr:              conf.Address(),
		Handler:           api,
		ReadTimeout:       conf.ReadTimeout.TimeDuration(),
		ReadHeaderTimeout: conf.ReadHeaderTimeout.TimeDuration(),
		WriteTimeout:      conf.WriteTimeout.TimeDuration(),
		IdleTimeout:       conf.IdleTimeout.TimeDuration(),
		MaxHeaderBytes:    conf.MaxHeaderBytes.Int(),
	}
	logger.Msg("start http server").Info(conf.Address())
	if err := api.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (api *Api) startMetricsServer() error {
	if len(api.conf.MetricsPort) == 0 {
		return nil
	}
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	api.metricsServer = &http.Server{
		Addr:    api.conf.MetricsAddress(),
		Handler: handler,
	}
	logger.Msg("start metrics http server").Info(api.conf.MetricsAddress())
	if err := api.metricsServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (api *Api) startDebugServer() error {
	if len(api.conf.DebugPort) == 0 {
		return nil
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if err := statsviz.Register(handler); err != nil {
		return fmt.Errorf("register handler for Statviz: %v", err)
	}
	api.debugServer = &http.Server{
		Addr:    api.conf.DebugAddress(),
		Handler: handler,
	}
	logger.Msg("start debug http server").Info(api.conf.DebugAddress())
	if err := api.debugServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (api *Api) ServeHTTP(rw http.ResponseWriter, httpReq *http.Request) {
	req := NewRequest(httpReq)
	res := NewResponse(rw, req.Id())

	logger.TraceId(req.Id()).Debug(newRequestLogStr(req.Method, req.Path()))
	defer logger.TraceId(req.Id()).Debug(finalRequestLogStr(req.Method, req.Path()))

	switch req.Path() {
	case "/v1/healthcheck":
		if req.Method != http.MethodGet {
			res.rw.WriteHeader(http.StatusOK)
			return
		}
	case "/v1/companies":
		api.companyHandlers(res, req)
	case "/v1/buildings":
		api.buildingHandlers(res, req)
	case "/v1/categories":
		api.categoriesHandlers(res, req)
	}
	if !res.err.IsEmpty() {
		logger.TraceId(req.Id()).Error(res.Error().String())
	}
	if !res.warn.IsEmpty() {
		logger.TraceId(req.Id()).Info(res.Warning().String())
	}
	res.writeHeaders()
	res.WriteBody()
}
