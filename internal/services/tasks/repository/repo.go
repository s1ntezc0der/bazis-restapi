package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/entity"
	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
)

type TaskRepository interface {
	Create(task *entity.Task) error
	GetByID(id int64) (*entity.Task, error)
	GetFiltered(filter *entity.TaskFilter) ([]entity.Task, error)
	Update(task *entity.Task) error
	AddHistory(history *entity.TaskHistory) error
	GetHistory(taskID int64) ([]entity.TaskHistory, error)
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(task *entity.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, assignee_id, team_id, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query, 
		task.Title, 
		task.Description, 
		task.Status, 
		task.AssigneeID, 
		task.TeamID, 
		task.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	task.ID = id

	return nil
}

func (r *taskRepo) GetByID(id int64) (*entity.Task, error) {
	query := `
		SELECT id, title, description, status, assignee_id, team_id, created_by, created_at, updated_at
		FROM tasks WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var task entity.Task
	err := row.Scan(
		&task.ID, 
		&task.Title, 
		&task.Description, 
		&task.Status, 
		&task.AssigneeID, 
		&task.TeamID,
		&task.CreatedBy, 
		&task.CreatedAt, 
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *taskRepo) GetFiltered(filter *entity.TaskFilter) ([]entity.Task, error) {
	query := `SELECT id, title, description, status, assignee_id, team_id, created_by, created_at, updated_at FROM tasks WHERE 1=1`
	args := []interface{}{}

	if filter.TeamID > 0 {
		query += " AND team_id = ?"
		args = append(args, filter.TeamID)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.AssigneeID > 0 {
		query += " AND assignee_id = ?"
		args = append(args, filter.AssigneeID)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID, 
			&task.Title, 
			&task.Description, 
			&task.Status, 
			&task.AssigneeID, 
			&task.TeamID, 
			&task.CreatedBy, 
			&task.CreatedAt, 
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepo) Update(task *entity.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, assignee_id = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query, 
		task.Title, 
		task.Description, 
		task.Status, 
		task.AssigneeID, 
		task.ID,
	)

	return err
}

func (r *taskRepo) AddHistory(history *entity.TaskHistory) error {
	query := `
		INSERT INTO task_history (task_id, changed_by, field, old_value, new_value)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query, 
		history.TaskID, 
		history.ChangedBy, 
		history.Field, 
		history.OldValue, 
		history.NewValue,
	)
	if err != nil {
		return fmt.Errorf("failed to add history: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	history.ID = id

	return nil
}

func (r *taskRepo) GetHistory(taskID int64) ([]entity.TaskHistory, error) {
	query := `
		SELECT id, task_id, changed_by, field, old_value, new_value, changed_at
		FROM task_history
		WHERE task_id = ?
		ORDER BY changed_at DESC
	`

	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer rows.Close()

	var history []entity.TaskHistory
	for rows.Next() {
		var h entity.TaskHistory
		if err := rows.Scan(
			&h.ID, 
			&h.TaskID, 
			&h.ChangedBy, 
			&h.Field, 
			&h.OldValue, 
			&h.NewValue, 
			&h.ChangedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}

func (r *taskRepo) GetTeamStats() ([]TeamStat, error) {
	query := `
		SELECT 
			t.id,
			t.name,
			COUNT(DISTINCT tm.user_id) AS members_count,
			COUNT(CASE WHEN tk.status = 'done' AND tk.updated_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 END) AS done_last_7_days
		FROM teams t
		LEFT JOIN team_members tm ON t.id = tm.team_id
		LEFT JOIN tasks tk ON t.id = tk.team_id
		GROUP BY t.id, t.name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []TeamStat
	for rows.Next() {
		var s TeamStat
		if err := rows.Scan(
			&s.TeamID, 
			&s.TeamName, 
			&s.MembersCount, 
			&s.DoneLast7Days,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func (r *taskRepo) GetTopCreators() ([]TopCreator, error) {
	query := `
		SELECT team_id, user_id, email, created_count, rank
		FROM (
			SELECT 
				t.team_id,
				u.id AS user_id,
				u.email,
				COUNT(t.id) AS created_count,
				ROW_NUMBER() OVER (PARTITION BY t.team_id ORDER BY COUNT(t.id) DESC) AS rank
			FROM tasks t
			JOIN users u ON t.created_by = u.id
			WHERE t.created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)
			GROUP BY t.team_id, u.id, u.email
		) ranked
		WHERE rank <= 3
		ORDER BY team_id, rank
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []TopCreator

	for rows.Next() {
		var c TopCreator
		if err := rows.Scan(
			&c.TeamID, 
			&c.UserID, 
			&c.Email, 
			&c.CreatedCount, 
			&c.Rank,
		); err != nil {
			return nil, err
		}
		creators = append(creators, c)
	}

	return creators, nil
}

func (r *taskRepo) GetInvalidAssignees() ([]InvalidAssignee, error) {
	query := `
		SELECT 
			t.id AS task_id,
			t.title,
			t.team_id,
			t.assignee_id,
			u.email AS assignee_email
		FROM tasks t
		JOIN users u ON t.assignee_id = u.id
		LEFT JOIN team_members tm ON t.team_id = tm.team_id AND t.assignee_id = tm.user_id
		WHERE t.assignee_id IS NOT NULL
			AND tm.user_id IS NULL
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invalid []
	
	for rows.Next() {
		var i InvalidAssignee
		if err := rows.Scan(
			&i.TaskID, 
			&i.TaskTitle, 
			&i.TeamID, 
			&i.AssigneeID, 
			&i.AssigneeEmail,
		); err != nil {
			return nil, err
		}
		invalid = append(invalid, i)
	}

	return invalid, nil
}

