package server

import (
	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/controller"
)

type Server struct {
	*HTTPServer
	*GrpcServer
}

func (s *Server) Run(config config.Server) {
	go s.HTTPServer.Run(config)
	go s.GrpcServer.Run(config)
}

func New(controller controller.Controller) *Server {
	return &Server{
		HTTPServer: newHTTPServer(controller),
		GrpcServer: newGrpcServer(controller),
	}
}
