package service

import (
	"strings"

	"github.com/mwozniak/missing-aed/internal/dto"
	"github.com/mwozniak/missing-aed/internal/model"
	"github.com/mwozniak/missing-aed/internal/repository"
)

const maxCommentLength = 200

type CommentService interface {
	Create(nodeID string, content string) (model.Comment, error)
	List(nodeID string) ([]model.Comment, error)
}

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) CommentService {
	return &commentService{commentRepository: commentRepository}
}

func (c commentService) Create(nodeID string, content string) (model.Comment, error) {
	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return model.Comment{}, dto.CommentNodeIDRequired
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return model.Comment{}, dto.CommentEmpty
	}
	if len([]rune(content)) > maxCommentLength {
		return model.Comment{}, dto.CommentTooLong
	}
	return c.commentRepository.Store(model.Comment{
		NodeID:  nodeID,
		Content: content,
	})
}

func (c commentService) List(nodeID string) ([]model.Comment, error) {
	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return nil, dto.CommentNodeIDRequired
	}
	return c.commentRepository.FindByNodeID(nodeID)
}
