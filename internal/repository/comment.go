package repository

import (
	"gorm.io/gorm"

	"github.com/mwozniak/missing-aed/internal/model"
)

type CommentRepository interface {
	Store(comment model.Comment) (model.Comment, error)
	FindByNodeID(nodeID string) ([]model.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func newCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (c commentRepository) Store(comment model.Comment) (model.Comment, error) {
	if err := c.db.Create(&comment).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (c commentRepository) FindByNodeID(nodeID string) ([]model.Comment, error) {
	var comments []model.Comment
	if err := c.db.Where("node_id = ?", nodeID).Order("created_at asc").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
