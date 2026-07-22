package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/comments/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/comments/usecase"
	"github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
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
func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
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

	var req entity.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.service.AddComment(taskID, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
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
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "id")

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	comments, err := h.service.GetComments(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

