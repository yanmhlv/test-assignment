package httpapi

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(handler *BookingHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", handler.CreateOrder)

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &Server{server: srv}
}

func (srv *Server) Start(addr string) error {
	srv.server.Addr = addr
	return srv.server.ListenAndServe()
}

func (srv *Server) Stop(ctx context.Context) error {
	return srv.server.Shutdown(ctx)
}
