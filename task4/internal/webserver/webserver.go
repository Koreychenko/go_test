package webserver

import (
	"context"
	"net/http"
)

type Webserver struct {
	mux *http.ServeMux
	srv *http.Server
}

func (srv *Webserver) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if srv.mux == nil {
		srv.mux = http.NewServeMux()
	}

	srv.mux.HandleFunc(pattern, handler)
}

func (srv *Webserver) Start() error {
	if srv.mux == nil {
		panic("there is no registered handlers")
	}

	srv.srv.Handler = srv.mux

	err := srv.srv.ListenAndServe()

	return err
}

func (srv *Webserver) Stop(ctx context.Context) error {
	err := srv.srv.Shutdown(ctx)

	return err
}

func NewWebserver(address string) *Webserver {
	return &Webserver{
		srv: &http.Server{Addr: address},
	}
}
