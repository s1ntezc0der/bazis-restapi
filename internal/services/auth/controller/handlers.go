package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mkk_bazis/internal/services/auth/entity"
	"mkk_bazis/internal/services/auth/usecase"
	"mkk_bazis/pkg/errors"
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
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		switch err {
		case errors.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
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
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		switch err {
		case errors.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
