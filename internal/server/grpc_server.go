package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/controller"
	pb "github.com/avran02/authentication/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcServer struct {
	pb.UnimplementedAuthServiceServer
	controller.Controller
}

func (s GrpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	slog.Info("Validating token")
	return s.Controller.ValidateToken(ctx, req)
}

func (s GrpcServer) Run(config config.Server) {
	serverEndpoint := fmt.Sprintf("%s:%s", config.Host, config.GRPCPort)
	slog.Info("Starting gRPC server on " + serverEndpoint)
	lis, err := net.Listen("tcp", serverEndpoint)
	if err != nil {
		slog.Error(fmt.Sprintf("can't listen on %s: \n%s", serverEndpoint, err.Error()))
		os.Exit(1)
	}

	slog.Info("Listening on " + serverEndpoint)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, s)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("authservice", grpc_health_v1.HealthCheckResponse_SERVING)

	slog.Info("Starting gRPC server")
	if err = grpcServer.Serve(lis); err != nil {
		slog.Error(fmt.Sprintf("can't start grpc server: \n%s", err.Error()))
		os.Exit(1)
	}
}

func newGrpcServer(controller controller.Controller) *GrpcServer {
	return &GrpcServer{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		Controller:                     controller,
	}
}
