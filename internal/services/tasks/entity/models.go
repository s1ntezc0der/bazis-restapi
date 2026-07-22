package entity

import "time"

type Task struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Status      string    `json:"status" db:"status"`
	AssigneeID  *int64    `json:"assignee_id,omitempty" db:"assignee_id"`
	TeamID      int64     `json:"team_id" db:"team_id"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TaskHistory struct {
	ID        int64     `json:"id" db:"id"`
	TaskID    int64     `json:"task_id" db:"task_id"`
	ChangedBy int64     `json:"changed_by" db:"changed_by"`
	Field     string    `json:"field" db:"field"`
	OldValue  *string   `json:"old_value" db:"old_value"`
	NewValue  *string   `json:"new_value" db:"new_value"`
	ChangedAt time.Time `json:"changed_at" db:"changed_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	AssigneeID  *int64 `json:"assignee_id"`
	TeamID      int64  `json:"team_id" validate:"required"`
}

type UpdateTaskRequest struct {
	Title      *string `json:"title"`
	Description *string `json:"description"`
	Status     *string `json:"status"`
	AssigneeID *int64  `json:"assignee_id"`
}

type TaskFilter struct {
	TeamID     int64
	Status     string
	AssigneeID int64
	Limit      int
	Offset     int
}

type TeamStat struct {
	TeamID         int64  `json:"team_id"`
	TeamName       string `json:"team_name"`
	MembersCount   int    `json:"members_count"`
	DoneLast7Days  int    `json:"done_last_7_days"`
}

type TopCreator struct {
	TeamID       int64  `json:"team_id"`
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	CreatedCount int    `json:"created_count"`
	Rank         int    `json:"rank"`
}

type InvalidAssignee struct {
	TaskID         int64  `json:"task_id"`
	TaskTitle      string `json:"task_title"`
	TeamID         int64  `json:"team_id"`
	AssigneeID     int64  `json:"assignee_id"`
	AssigneeEmail  string `json:"assignee_email"`
}

