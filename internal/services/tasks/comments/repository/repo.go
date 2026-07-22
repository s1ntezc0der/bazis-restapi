package repository

import (
	"database/sql"
	"fmt"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/comments/entity"
	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
)

type CommentRepository interface {
	Create(comment *entity.TaskComment) error
	GetByTaskID(taskID int64) ([]entity.TaskComment, error)
}

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(comment *entity.TaskComment) error {
	query := `
		INSERT INTO task_comments (task_id, user_id, content)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
        query, 
        comment.TaskID, 
        comment.UserID, 
        comment.Content,
    )
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	comment.ID = id

	return nil
}

func (r *commentRepo) GetByTaskID(taskID int64) ([]entity.TaskComment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at
		FROM task_comments
		WHERE task_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []entity.TaskComment

	for rows.Next() {
		var c entity.TaskComment
		if err := rows.Scan(
            &c.ID, 
            &c.TaskID, 
            &c.UserID, 
            &c.Content, 
            &c.CreatedAt,
        ); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

