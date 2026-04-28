package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/mwozniak/missing-aed/internal/dto"
	"github.com/mwozniak/missing-aed/internal/service"
)

type CommentController interface {
	Index(c *echo.Context) error
	Store(c *echo.Context) error
}

type commentController struct {
	commentService service.CommentService
}

func newCommentController(commentService service.CommentService) CommentController {
	return &commentController{commentService: commentService}
}

func (cc commentController) Index(c *echo.Context) error {
	nodeID := c.Param("nodeID")
	comments, err := cc.commentService.List(nodeID)
	if err != nil {
		return mapCommentError(err)
	}
	return c.JSON(http.StatusOK, comments)
}

func (cc commentController) Store(c *echo.Context) error {
	nodeID := c.Param("nodeID")
	var request dto.CreateCommentRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	comment, err := cc.commentService.Create(nodeID, request.Content)
	if err != nil {
		return mapCommentError(err)
	}
	return c.JSON(http.StatusCreated, comment)
}

func mapCommentError(err error) error {
	if errors.Is(err, dto.CommentEmpty) || errors.Is(err, dto.CommentTooLong) || errors.Is(err, dto.CommentNodeIDRequired) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return err
}
