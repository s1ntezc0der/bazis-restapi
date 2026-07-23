package auth

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "mkk_bazis/internal/services/auth/entity"
    "mkk_bazis/internal/services/auth/usecase"
    "mkk_bazis/pkg/jwt"
)

type MockAuthRepo struct {
    mock.Mock
}

func (m *MockAuthRepo) CreateUser(user *entity.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockAuthRepo) GetUserByEmail(email string) (*entity.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthRepo) GetUserByID(id int64) (*entity.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
    mockRepo := new(MockAuthRepo)
    jwtConfig := jwt.NewJWT("secret", 24*time.Hour)
    service := usecase.NewAuthService(mockRepo, jwtConfig)

    req := &entity.RegisterRequest{
        Email:    "test@example.com",
        Password: "123456",
        Name:     "Test User",
    }

    mockRepo.On("GetUserByEmail", req.Email).Return(nil, nil)
    mockRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

    user, err := service.Register(req)

    assert.NoError(t, err)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, req.Name, user.Name)
    assert.NotEmpty(t, user.PasswordHash)
    mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
    mockRepo := new(MockAuthRepo)
    jwtConfig := jwt.NewJWT("secret", 24*time.Hour)
    service := usecase.NewAuthService(mockRepo, jwtConfig)

    existingUser := &entity.User{Email: "test@example.com"}
    req := &entity.RegisterRequest{
        Email:    "test@example.com",
        Password: "123456",
        Name:     "Test User",
    }

    mockRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil)

    user, err := service.Register(req)

    assert.Nil(t, user)
    assert.Error(t, err)
    assert.Equal(t, "email already exists", err.Error())
    mockRepo.AssertExpectations(t)
}

