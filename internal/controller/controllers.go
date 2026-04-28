package controller

import (
	"github.com/labstack/echo/v5"

	"github.com/mwozniak/missing-aed/internal/service"
)

type Controllers interface {
	Route(e *echo.Echo)
}

type controllers struct {
	aedController AedController
}

func NewControllers(services service.Services) Controllers {
	return &controllers{
		aedController: newAedController(services.Matcher()),
	}
}

func (c controllers) Route(e *echo.Echo) {
	e.GET("/api/missing-aeds", c.aedController.Missing)
}
