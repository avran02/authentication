package controller

import (
	"github.com/avran02/authentication/internal/config"
	"github.com/avran02/authentication/internal/service"
)

type Controller interface {
	HTTPController
	GrpcController
}

type controller struct {
	HTTPController
	GrpcController
}

func New(service service.Service, cookieConfig config.CookieConfig) Controller {
	return &controller{
		HTTPController: newHTTPController(service, cookieConfig),
		GrpcController: newGrpcController(service),
	}
}
