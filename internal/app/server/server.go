package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log        *slog.Logger
	httpServer *http.Server
	handler    http.Handler
	port       string
}

func New(log *slog.Logger, port string, handler http.Handler) *Server {
	return &Server{
		log:     log,
		handler: handler,
		port:    port,
	}
}

func (s *Server) MustRun() {
	err := s.Run()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "server.Run"
	s.log.With(slog.String("op", op)).
		Info("starting server")

	s.httpServer = &http.Server{
		Addr:           ":" + s.port,
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	const op = "server.Stop"
	s.log.With(slog.String("op", op)).Info("stopping server")

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.log.Error("failed to shutdown http server", "error", err)
	}
}
