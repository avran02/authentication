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
	"github.com/go-chi/cors"
	swagger "github.com/swaggo/http-swagger"
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

func newHTTPServer(controller controller.Controller, debug bool) *HTTPServer {
	s := &HTTPServer{
		controller: controller,
	}
	corsOpts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            debug,
	}

	main := chi.NewMux()
	main.Use(middleware.Logger)
	main.Use(middleware.Recoverer)
	main.Use(cors.Handler(corsOpts))

	main.Get("/docs/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.yml")
	})
	main.Get("/swagger/*", swagger.Handler(
		swagger.URL("/docs/openapi.yml"),
	))
	main.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})

	router := s.routes()
	main.Mount("/api/v1", router)
	s.router = main

	return s
}
