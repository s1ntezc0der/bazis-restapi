package teams

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "mkk_bazis/internal/services/teams/entity"
    "mkk_bazis/internal/services/teams/usecase"
)

type MockTeamRepo struct {
    mock.Mock
}

func (m *MockTeamRepo) Create(team *entity.Team) error {
    args := m.Called(team)
    return args.Error(0)
}

func (m *MockTeamRepo) GetByID(id int64) (*entity.Team, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepo) GetByUserID(userID int64) ([]entity.Team, error) {
    args := m.Called(userID)
    return args.Get(0).([]entity.Team), args.Error(1)
}

func (m *MockTeamRepo) AddMember(member *entity.TeamMember) error {
    args := m.Called(member)
    return args.Error(0)
}

func (m *MockTeamRepo) GetMember(teamID, userID int64) (*entity.TeamMember, error) {
    args := m.Called(teamID, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.TeamMember), args.Error(1)
}

func (m *MockTeamRepo) GetMembers(teamID int64) ([]entity.TeamMember, error) {
    args := m.Called(teamID)
    return args.Get(0).([]entity.TeamMember), args.Error(1)
}

func (m *MockTeamRepo) UpdateRole(teamID, userID int64, role string) error {
    args := m.Called(teamID, userID, role)
    return args.Error(0)
}

type MockAuthRepo struct {
    mock.Mock
}

func (m *MockAuthRepo) CreateUser(user *entity.User) error {
    return nil
}
func (m *MockAuthRepo) GetUserByEmail(email string) (*entity.User, error) {
    return nil, nil
}
func (m *MockAuthRepo) GetUserByID(id int64) (*entity.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateTeam_Success(t *testing.T) {
    mockTeamRepo := new(MockTeamRepo)
    mockAuthRepo := new(MockAuthRepo)
    service := usecase.NewTeamService(mockTeamRepo, mockAuthRepo)

    req := &entity.CreateTeamRequest{
        Name:        "Team Alpha",
        Description: "Test team",
    }
    userID := int64(1)

    mockTeamRepo.On("Create", mock.AnythingOfType("*entity.Team")).Return(nil)
    mockTeamRepo.On("AddMember", mock.AnythingOfType("*entity.TeamMember")).Return(nil)

    team, err := service.CreateTeam(userID, req)

    assert.NoError(t, err)
    assert.Equal(t, req.Name, team.Name)
    assert.Equal(t, req.Description, team.Description)
    assert.Equal(t, userID, team.CreatedBy)
    mockTeamRepo.AssertExpectations(t)
}

