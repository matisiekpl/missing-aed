package repository

import "gorm.io/gorm"

type Repositories interface {
	Comment() CommentRepository
}

type repositories struct {
	commentRepository CommentRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	return &repositories{
		commentRepository: newCommentRepository(db),
	}
}

func (r *repositories) Comment() CommentRepository {
	return r.commentRepository
}
