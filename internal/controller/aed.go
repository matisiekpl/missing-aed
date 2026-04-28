package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/mwozniak/missing-aed/internal/dto"
	"github.com/mwozniak/missing-aed/internal/service"
)

type AedController interface {
	Missing(c *echo.Context) error
}

type aedController struct {
	matcherService service.MatcherService
}

func newAedController(matcherService service.MatcherService) AedController {
	return &aedController{matcherService: matcherService}
}

func (a aedController) Missing(c *echo.Context) error {
	missing, err := a.matcherService.Missing()
	if err != nil {
		if errors.Is(err, dto.OsmNotReady) {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, missing)
}
