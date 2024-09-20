package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/controller"
	"github.com/avran02/authentication/internal/pkg/jwt"
	"github.com/avran02/authentication/internal/repo"
	"github.com/avran02/authentication/internal/server"
	"github.com/avran02/authentication/internal/service"
	"github.com/avran02/authentication/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	server     *server.Server
	config     *config.Config
	controller controller.Controller
}

func (app *App) Run() {
	host := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
	slog.Info("Starting gRPC server on " + host)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		slog.Error(fmt.Sprintf("can't listen on %s: \n%s", host, err.Error()))
		os.Exit(1)
	}

	slog.Info("Listening on " + host)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, app.server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("authservice", grpc_health_v1.HealthCheckResponse_SERVING)

	slog.Info("Starting gRPC server")
	if err = grpcServer.Serve(lis); err != nil {
		slog.Error(fmt.Sprintf("can't start grpc server: \n%s", err.Error()))
		os.Exit(1)
	}
}

func New() *App {
	config := config.New()

	repo := repo.New(&config.DB)
	JWTGenerator := jwt.NewJwtGenerator(config.JWT)
	service := service.New(repo, JWTGenerator)
	controller := controller.New(service)
	server := server.New(controller)

	return &App{
		config:     config,
		controller: controller,
		server:     server,
	}
}
