package usecase

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"mkk_bazis/internal/services/auth/entity"
	"mkk_bazis/internal/services/auth/repository"
	"mkk_bazis/pkg/errors"
	"mkk_bazis/pkg/jwt"
)

type AuthService interface {
	Register(req *entity.RegisterRequest) (*entity.User, error)
	Login(req *entity.LoginRequest) (*entity.LoginResponse, error)
	GetUserByID(id int64) (*entity.User, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtConfig *jwt.JWTConfig
}

func NewAuthService(repo repository.AuthRepository, jwtConfig *jwt.JWTConfig) AuthService {
	return &authService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *authService) Register(req *entity.RegisterRequest) (*entity.User, error) {
	existing, _ := s.repo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	token, err := s.jwtConfig.Generate(user.ID, user.Email)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &entity.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GetUserByID(id int64) (*entity.User, error) {
	return s.repo.GetUserByID(id)
}

