package repository

import (
	"database/sql"
	"fmt"

	"mkk_bazis/internal/services/auth/entity"
	"mkk_bazis/pkg/errors"
)

type AuthRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id int64) (*entity.User, error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(user *entity.User) error {
	query := `
		INSERT INTO users (email, password_hash, name)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
		query, 
		user.Email, 
		user.PasswordHash, 
		user.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	
	return nil
}

func (r *authRepo) GetUserByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	row := r.db.QueryRow(query, email)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *authRepo) GetUserByID(id int64) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

