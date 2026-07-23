package repository

import (
	"database/sql"
	"fmt"

	"mkk_bazis/internal/services/teams/entity"
	"mkk_bazis/pkg/errors"
)

type TeamRepository interface {
	Create(team *entity.Team) error
	GetByID(id int64) (*entity.Team, error)
	GetByUserID(userID int64) ([]entity.Team, error)
	AddMember(member *entity.TeamMember) error
	GetMember(teamID, userID int64) (*entity.TeamMember, error)
	GetMembers(teamID int64) ([]entity.TeamMember, error)
	UpdateRole(teamID, userID int64, role string) error
}

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(team *entity.Team) error {
	query := `
		INSERT INTO teams (name, description, created_by)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
        query, 
        team.Name, 
        team.Description, 
        team.CreatedBy,
    )
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	team.ID = id

	return nil
}

func (r *teamRepo) GetByID(id int64) (*entity.Team, error) {
	query := `
		SELECT id, name, description, created_by, created_at, updated_at
		FROM teams WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var team entity.Team
	err := row.Scan(
        &team.ID, 
        &team.Name, 
        &team.Description, 
        &team.CreatedBy, 
        &team.CreatedAt, 
        &team.UpdatedAt,
    )
	if err == sql.ErrNoRows {
		return nil, errors.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

func (r *teamRepo) GetByUserID(userID int64) ([]entity.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_by, t.created_at, t.updated_at
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	defer rows.Close()

	var teams []entity.Team

	for rows.Next() {
		var team entity.Team
		if err := rows.Scan(
            &team.ID, 
            &team.Name, 
            &team.Description, 
            &team.CreatedBy, 
            &team.CreatedAt, 
            &team.UpdatedAt,
        ); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *teamRepo) AddMember(member *entity.TeamMember) error {
	query := `
		INSERT INTO team_members (user_id, team_id, role)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
        query, 
        member.UserID, 
        member.TeamID, 
        member.Role,
    )
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	member.ID = id
	
    return nil
}

func (r *teamRepo) GetMember(teamID, userID int64) (*entity.TeamMember, error) {
	query := `
		SELECT id, user_id, team_id, role, joined_at
		FROM team_members
		WHERE team_id = ? AND user_id = ?
	`

	row := r.db.QueryRow(query, teamID, userID)

	var member entity.TeamMember
	
    err := row.Scan(
        &member.ID, 
        &member.UserID, 
        &member.TeamID, 
        &member.Role, 
        &member.JoinedAt,
    )
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	
    return &member, nil
}

func (r *teamRepo) GetMembers(teamID int64) ([]entity.TeamMember, error) {
	query := `
		SELECT id, user_id, team_id, role, joined_at
		FROM team_members
		WHERE team_id = ?
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	var members []entity.TeamMember
	
    for rows.Next() {
		var m entity.TeamMember
		if err := rows.Scan(
            &m.ID, 
            &m.UserID, 
            &m.TeamID, 
            &m.Role, &
            m.JoinedAt,
        ); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, nil
}

func (r *teamRepo) UpdateRole(teamID, userID int64, role string) error {
	query := `
        UPDATE team_members SET role = ? 
        WHERE team_id = ? AND user_id = ?
    `
    
	_, err := r.db.Exec(
        query, 
        role, 
        teamID, 
        userID,
    )

	return err
}

