package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mkk_bazis/internal/services/tasks/entity"
	"mkk_bazis/internal/services/tasks/usecase"
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
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req entity.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	task, err := h.service.CreateTask(userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
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
func (h *TaskHandler) GetTasks(c *gin.Context) {
	filter := &entity.TaskFilter{
		Limit:  10,
		Offset: 0,
	}

	if teamID, err := strconv.ParseInt(c.Query("team_id"), 10, 64); err == nil {
		filter.TeamID = teamID
	}
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}
	if assigneeID, err := strconv.ParseInt(c.Query("assignee_id"), 10, 64); err == nil {
		filter.AssigneeID = assigneeID
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.Atoi(c.Query("offset")); err == nil {
		filter.Offset = offset
	}

	tasks, err := h.service.GetTasks(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
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
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req entity.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	task, err := h.service.UpdateTask(userID.(int64), taskID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
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
func (h *TaskHandler) GetHistory(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	history, err := h.service.GetHistory(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

