package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	controller controller.Controller
	router     *chi.Mux
}

func (s *HTTPServer) routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/register", s.controller.Register)
	r.Post("/login", s.controller.Login)
	r.Post("/refresh-tokens", s.controller.RefreshTokens)
	r.Post("/logout", s.controller.Logout)

	return r
}

func (s *HTTPServer) Run(config config.Server) {
	serverEndpoint := fmt.Sprintf("%s:%s", config.Host, config.HTTPPort)
	slog.Info("Starting http server at " + serverEndpoint)
	server := http.Server{ //nolint:gosec
		Addr:    serverEndpoint,
		Handler: s.router,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("can't statr http server", "error", err.Error())
		os.Exit(1)
	}
}

func newHTTPServer(controller controller.Controller) *HTTPServer {
	s := &HTTPServer{
		controller: controller,
	}
	main := chi.NewMux()
	main.Use(middleware.Logger)
	main.Use(middleware.Recoverer)

	router := s.routes()
	main.Mount("/", router)
	s.router = main

	return s
}
