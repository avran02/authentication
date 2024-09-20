package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avran02/authentication/internal/service"
	pb "github.com/avran02/authentication/pb"
)

type Controller interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	RefreshTokens(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
}

// implements pb.AuthServiceServer.
type controller struct {
	servcie service.Service
}

func (c *controller) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := c.servcie.Register(ctx, req.Email, req.Username, req.Password); err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	resp := &pb.RegisterResponse{
		Success: true,
	}

	return resp, nil
}

func (c *controller) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, refreshToken, err := c.servcie.Login(ctx, req.Username, req.Password)
	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("failed to login user: %w", err)
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *controller) RefreshTokens(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := c.servcie.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *controller) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	id, err := c.servcie.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	return &pb.ValidateTokenResponse{
		Id: id,
	}, nil
}

func (c *controller) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	ok, err := c.servcie.Logout(ctx, req.AccessToken)
	if err != nil {
		err = fmt.Errorf("failed to logout: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &pb.LogoutResponse{
		Success: ok,
	}, nil
}

func New(service service.Service) Controller {
	return &controller{
		servcie: service,
	}
}
