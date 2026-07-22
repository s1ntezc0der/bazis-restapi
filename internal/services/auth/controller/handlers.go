package controller

import (
	"encoding/json"
	"net/http"

	"github.com/s1ntezc0der/bazis-restapi/internal/services/auth/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/auth/usecase"
	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
)

type AuthHandler struct {
	service usecase.AuthService
}

func NewAuthHandler(service usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя с email, паролем и именем
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} entity.User "Пользователь создан"
// @Failure 400 {string} string "Invalid request body"
// @Failure 409 {string} string "Email already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		switch err {
		case errors.ErrEmailAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Данные для входа"
// @Success 200 {object} entity.LoginResponse "Токен и данные пользователя"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		switch err {
		case errors.ErrInvalidCredentials:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

