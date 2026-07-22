package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/repository"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/teams/repository"
	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
)

type TaskService interface {
	CreateTask(userID int64, req *entity.CreateTaskRequest) (*entity.Task, error)
	GetTasks(filter *entity.TaskFilter) ([]entity.Task, error)
	UpdateTask(userID int64, taskID int64, req *entity.UpdateTaskRequest) (*entity.Task, error)
	GetHistory(taskID int64) ([]entity.TaskHistory, error)
}

type taskService struct {
    taskRepo repository.TaskRepository
    teamRepo teamsRepo.TeamRepository
    cache    *cache.Cache
}

func NewTaskService(taskRepo repository.TaskRepository, teamRepo teamsRepo.TeamRepository, cache *cache.Cache) TaskService {
	return &taskService{
        taskRepo: taskRepo,
        teamRepo: teamRepo,
        cache:    cache,
    }
}

func (s *taskService) CreateTask(userID int64, req *entity.CreateTaskRequest) (*entity.Task, error) {
	member, err := s.teamRepo.GetMember(req.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.ErrNotMember
	}

	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      "todo",
		AssigneeID:  req.AssigneeID,
		TeamID:      req.TeamID,
		CreatedBy:   userID,
	}

	if req.AssigneeID != nil {
		member, err := s.teamRepo.GetMember(req.TeamID, *req.AssigneeID)
		if err != nil {
			return nil, err
		}
		if member == nil {
			return nil, errors.ErrAssigneeNotInTeam
		}
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTasks(filter *entity.TaskFilter) ([]entity.Task, error) {
	return s.taskRepo.GetFiltered(filter)
}

func (s *taskService) UpdateTask(userID int64, taskID int64, req *entity.UpdateTaskRequest) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	member, err := s.teamRepo.GetMember(task.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.ErrNotMember
	}
	if member.Role != "owner" && member.Role != "admin" && task.CreatedBy != userID {
		return nil, errors.ErrForbidden
	}

	oldTask := *task

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.AssigneeID != nil {
		if *req.AssigneeID != 0 {
			member, err := s.teamRepo.GetMember(task.TeamID, *req.AssigneeID)
			if err != nil {
				return nil, err
			}
			if member == nil {
				return nil, errors.ErrAssigneeNotInTeam
			}
		}
		task.AssigneeID = req.AssigneeID
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	s.saveHistory(&oldTask, task, userID)

	return task, nil
}

func (s *taskService) GetHistory(taskID int64) ([]entity.TaskHistory, error) {
	return s.taskRepo.GetHistory(taskID)
}

func (s *taskService) saveHistory(old, new *entity.Task, userID int64) {
	fields := map[string]interface{}{
		"title":       struct{ Old, New string }{Old: old.Title, New: new.Title},
		"description": struct{ Old, New string }{Old: old.Description, New: new.Description},
		"status":      struct{ Old, New string }{Old: old.Status, New: new.Status},
		"assignee_id": struct{ Old, New interface{} }{Old: old.AssigneeID, New: new.AssigneeID},
	}

	for field, val := range fields {
		oldVal := reflect.ValueOf(val).FieldByName("Old").Interface()
		newVal := reflect.ValueOf(val).FieldByName("New").Interface()

		if !reflect.DeepEqual(oldVal, newVal) {
			history := &entity.TaskHistory{
				TaskID:    new.ID,
				ChangedBy: userID,
				Field:     field,
			}

			if oldVal != nil {
				oldStr := ""
				switch v := oldVal.(type) {
				case string:
					oldStr = v
				case *int64:
					if v != nil {
						oldStr = fmt.Sprintf("%d", *v)
					}
				}
				history.OldValue = &oldStr
			}

			if newVal != nil {
				newStr := ""
				switch v := newVal.(type) {
				case string:
					newStr = v
				case *int64:
					if v != nil {
						newStr = fmt.Sprintf("%d", *v)
					}
				}
				history.NewValue = &newStr
			}

			s.taskRepo.AddHistory(history)
		}
	}
}

func (s *taskService) GetTasks(filter *entity.TaskFilter) ([]entity.Task, error) {
    cacheKey := fmt.Sprintf("tasks:team:%d:status:%s:assignee:%d", filter.TeamID, filter.Status, filter.AssigneeID)
    
    var tasks []entity.Task
    if err := s.cache.Get(context.Background(), cacheKey, &tasks); err == nil {
        return tasks, nil
    }

    tasks, err := s.taskRepo.GetFiltered(filter)
    if err != nil {
        return nil, err
    }

    s.cache.Set(context.Background(), cacheKey, tasks, 5*time.Minute)
    return tasks, nil
}

