package tasks

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "mkk_bazis/internal/services/tasks/entity"
    "mkk_bazis/internal/services/tasks/usecase"
)

type MockTaskRepo struct {
    mock.Mock
}

func (m *MockTaskRepo) Create(task *entity.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *MockTaskRepo) GetByID(id int64) (*entity.Task, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskRepo) GetFiltered(filter *entity.TaskFilter) ([]entity.Task, error) {
    args := m.Called(filter)
    return args.Get(0).([]entity.Task), args.Error(1)
}

func (m *MockTaskRepo) Update(task *entity.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *MockTaskRepo) AddHistory(history *entity.TaskHistory) error {
    args := m.Called(history)
    return args.Error(0)
}

func (m *MockTaskRepo) GetHistory(taskID int64) ([]entity.TaskHistory, error) {
    args := m.Called(taskID)
    return args.Get(0).([]entity.TaskHistory), args.Error(1)
}

type MockTeamRepo struct {
    mock.Mock
}

func (m *MockTeamRepo) GetMember(teamID, userID int64) (*entity.TeamMember, error) {
    args := m.Called(teamID, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.TeamMember), args.Error(1)
}

func TestCreateTask_Success(t *testing.T) {
    mockTaskRepo := new(MockTaskRepo)
    mockTeamRepo := new(MockTeamRepo)
    service := usecase.NewTaskService(mockTaskRepo, mockTeamRepo)

    req := &entity.CreateTaskRequest{
        Title:  "Test Task",
        TeamID: 1,
    }
    userID := int64(1)

    mockTeamRepo.On("GetMember", req.TeamID, userID).Return(&entity.TeamMember{
		ID: 1, 
		UserID: userID, 
		TeamID: req.TeamID, 
		Role: "member",
	}, nil)
    mockTaskRepo.On("Create", mock.AnythingOfType("*entity.Task")).Return(nil)

    task, err := service.CreateTask(userID, req)

    assert.NoError(t, err)
    assert.Equal(t, req.Title, task.Title)
    assert.Equal(t, req.TeamID, task.TeamID)
    assert.Equal(t, userID, task.CreatedBy)
    assert.Equal(t, "todo", task.Status)
    mockTaskRepo.AssertExpectations(t)
}
