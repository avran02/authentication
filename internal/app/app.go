package app

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/controller"
	"github.com/avran02/authentication/internal/pkg/jwt"
	"github.com/avran02/authentication/internal/repo"
	"github.com/avran02/authentication/internal/server"
	"github.com/avran02/authentication/internal/service"
	"github.com/avran02/authentication/logger"
)

type App struct {
	server     *server.Server
	config     *config.Config
	controller controller.Controller
}

func (app *App) Run() {
	app.server.Run(app.config.Server)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signals
	slog.Info("shutdown server", "signal", sig.String())
	os.Exit(0)
}

func New() *App {
	config := config.New()
	logger.Setup(config.Server)
	slog.Debug("config loaded", "config", config)
	debug := false
	if strings.ToLower(config.Server.LogLevel) == "debug" {
		debug = true
	}

	repo := repo.New(&config.DB)
	JWTGenerator := jwt.NewJwtGenerator(config.JWT)
	service := service.New(repo, JWTGenerator)
	controller := controller.New(service, config.Cookie)
	server := server.New(controller, debug, config.CORS)

	return &App{
		config:     config,
		controller: controller,
		server:     server,
	}
}
