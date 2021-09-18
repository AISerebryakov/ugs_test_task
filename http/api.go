package http

import "net/http"

type handler struct {
	Api
}

type Api struct {
	server  *http.Server
	conf    Config
	firmMng FirmManager
}

func NewApi(conf Config) (api Api) {
	api.conf = conf
	api.firmMng = conf.FirmManager
	return api
}

//todo: add metrics server
//todo: add debug server

func (api *Api) startServer() error {
	conf := api.conf
	api.server = &http.Server{
		Addr:              conf.Address(),
		Handler:           handler{*api},
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	if err := api.server.ListenAndServe(); err != nil {
		//todo: handle error
		return err
	}
	return nil
}

func (h handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case firmPath:
		h.firmHandlers(rw, req)
	}
}
