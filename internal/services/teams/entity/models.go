package entity

import "time"

type Team struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TeamMember struct {
	ID       int64     `json:"id" db:"id"`
	UserID   int64     `json:"user_id" db:"user_id"`
	TeamID   int64     `json:"team_id" db:"team_id"`
	Role     string    `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type InviteRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type TeamWithMembers struct {
	Team
	Members []TeamMember `json:"members"`
}

