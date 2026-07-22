package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/usecase"
	"github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
)

type TaskHandler struct {
	service usecase.TaskService
}

func NewTaskHandler(service usecase.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary Создать задачу
// @Description Создаёт новую задачу в команде (только член команды)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreateTaskRequest true "Данные задачи"
// @Success 201 {object} entity.Task "Задача создана"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req entity.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTasks godoc
// @Summary Список задач с фильтрацией
// @Description Возвращает задачи с фильтрацией по команде, статусу, исполнителю и пагинацией
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param team_id query int false "ID команды"
// @Param status query string false "Статус задачи (todo, in_progress, done)"
// @Param assignee_id query int false "ID исполнителя"
// @Param limit query int false "Лимит записей" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} entity.Task "Список задач"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks [get]
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := &entity.TaskFilter{
		Limit:  10,
		Offset: 0,
	}

	if teamID, err := strconv.ParseInt(query.Get("team_id"), 10, 64); err == nil {
		filter.TeamID = teamID
	}
	if status := query.Get("status"); status != "" {
		filter.Status = status
	}
	if assigneeID, err := strconv.ParseInt(query.Get("assignee_id"), 10, 64); err == nil {
		filter.AssigneeID = assigneeID
	}
	if limit, err := strconv.Atoi(query.Get("limit")); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.Atoi(query.Get("offset")); err == nil {
		filter.Offset = offset
	}

	tasks, err := h.service.GetTasks(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// UpdateTask godoc
// @Summary Обновить задачу
// @Description Обновляет поля задачи (только owner/admin или создатель)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param request body entity.UpdateTaskRequest true "Данные для обновления"
// @Success 200 {object} entity.Task "Задача обновлена"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req entity.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.UpdateTask(userID, taskID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// GetHistory godoc
// @Summary История изменений задачи
// @Description Возвращает историю изменений задачи
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {array} entity.TaskHistory "История изменений"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks/{id}/history [get]
func (h *TaskHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetHistory(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

