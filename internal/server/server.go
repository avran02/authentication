package server

import (
	"context"
	"log/slog"

	"github.com/avran02/authentication/internal/controller"
	pb "github.com/avran02/authentication/pb"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	controller.Controller
}

func (s Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	slog.Info("Registering user")
	return s.Controller.Register(ctx, req)
}

func (s Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	slog.Info("Logging in user")
	return s.Controller.Login(ctx, req)
}

func (s Server) RefreshTokens(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	slog.Info("Refreshing token")
	return s.Controller.RefreshTokens(ctx, req)
}

func (s Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	slog.Info("Validating token")
	return s.Controller.ValidateToken(ctx, req)
}

func (s Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	slog.Info("Logging out user")
	return s.Controller.Logout(ctx, req)
}

func New(controller controller.Controller) *Server {
	return &Server{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		Controller:                     controller,
	}
}
