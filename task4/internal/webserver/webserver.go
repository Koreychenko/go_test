package webserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const timeout = 2 * time.Second

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
		return errors.New("there is no registered handlers")
	}

	srv.srv.Handler = srv.mux

	return srv.srv.ListenAndServe()
}

func (srv *Webserver) Stop(ctx context.Context) error {
	return srv.srv.Shutdown(ctx)
}

func NewWebserver(address string) *Webserver {
	return &Webserver{
		srv: &http.Server{Addr: address},
	}
}

func WrapHandler(handler func(*http.Request) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(request.Context(), timeout)
		defer cancel()

		request = request.WithContext(ctx)

		bytes, err := handler(request)

		if err == nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			if len(bytes) > 0 {
				writer.Write(bytes)
			}

			return
		}

		switch {
		case errors.Is(err, &ValidationError{}):
			http.Error(writer, err.Error(), http.StatusBadRequest)
		case errors.Is(err, &InvalidRequestError{}):
			http.Error(writer, err.Error(), http.StatusBadRequest)
		case errors.Is(err, &NotFoundError{}):
			http.Error(writer, err.Error(), http.StatusNotFound)
		default:
			http.Error(writer, "internal error", http.StatusInternalServerError)
		}
	}
}
