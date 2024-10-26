package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avran02/authentication/internal/service"
	pb "github.com/avran02/authentication/pb"
)

type GrpcController interface {
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
}

// implements pb.AuthServiceServer.
type grpcController struct {
	service service.Service
}

func (c *grpcController) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	id, err := c.service.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	return &pb.ValidateTokenResponse{
		Id: id,
	}, nil
}

func newGrpcController(service service.Service) GrpcController {
	return &grpcController{
		service: service,
	}
}
