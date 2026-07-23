package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mkk_bazis/internal/services/tasks/comments/entity"
	"mkk_bazis/internal/services/tasks/comments/usecase"
)

type CommentHandler struct {
	service usecase.CommentService
}

func NewCommentHandler(service usecase.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// AddComment godoc
// @Summary Добавить комментарий к задаче
// @Description Добавляет комментарий к задаче
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param request body entity.CreateCommentRequest true "Текст комментария"
// @Success 201 {object} entity.TaskComment "Комментарий создан"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks/{id}/comments [post]
func (h *CommentHandler) AddComment(c *gin.Context) {
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

	var req entity.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	comment, err := h.service.AddComment(taskID, userID.(int64), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments godoc
// @Summary Получить все комментарии к задаче
// @Description Возвращает все комментарии к задаче
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {array} entity.TaskComment "Список комментариев"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/tasks/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	comments, err := h.service.GetComments(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
