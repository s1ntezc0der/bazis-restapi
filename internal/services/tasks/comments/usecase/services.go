package usecase

import (
	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/comments/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/comments/repository"
)

type CommentService interface {
	AddComment(taskID, userID int64, content string) (*entity.TaskComment, error)
	GetComments(taskID int64) ([]entity.TaskComment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) AddComment(taskID, userID int64, content string) (*entity.TaskComment, error) {
	comment := &entity.TaskComment{
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetComments(taskID int64) ([]entity.TaskComment, error) {
	return s.repo.GetByTaskID(taskID)
}

