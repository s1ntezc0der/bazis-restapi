package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mkk_bazis/internal/services/teams/entity"
	"mkk_bazis/internal/services/teams/usecase"
)

type TeamHandler struct {
	service usecase.TeamService
}

func NewTeamHandler(service usecase.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// CreateTeam godoc
// @Summary Создать команду
// @Description Создаёт новую команду и назначает текущего пользователя владельцем
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreateTeamRequest true "Данные команды"
// @Success 201 {object} entity.Team "Команда создана"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req entity.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	team, err := h.service.CreateTeam(userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetUserTeams godoc
// @Summary Список команд пользователя
// @Description Возвращает все команды, в которых состоит текущий пользователь
// @Tags teams
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Team "Список команд"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/teams [get]
func (h *TeamHandler) GetUserTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teams, err := h.service.GetUserTeams(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// InviteUser godoc
// @Summary Пригласить пользователя в команду
// @Description Приглашает пользователя в команду (только owner/admin)
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID команды"
// @Param request body entity.InviteRequest true "ID приглашаемого пользователя"
// @Success 200 {string} string "Invitation sent"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Team or user not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/teams/{id}/invite [post]
func (h *TeamHandler) InviteUser(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	inviterID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req entity.InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.InviteUser(teamID, inviterID.(int64), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation sent"})
}
