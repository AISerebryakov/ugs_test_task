package http

import "net/http"

type handler struct {
	*Api
}

type Api struct {
	server *http.Server
	conf   Config
}

func NewApi(conf Config) (api Api) {

	return api
}

//todo: add metrics server
//todo: add debug server

func (api *Api) startServer() error {
	conf := api.conf
	api.server = &http.Server{
		Addr:              conf.Address(),
		Handler:           handler{api},
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
	case "/api/v1/firm":
		h.firmHandlers(rw, req)
	}
}
