package controller

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "github.com/s1ntezc0der/bazis-restapi/internal/services/teams/entity"
    "github.com/s1ntezc0der/bazis-restapi/internal/services/teams/usecase"
    "github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
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
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var req entity.CreateTeamRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    team, err := h.service.CreateTeam(userID, &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(team)
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
func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    teams, err := h.service.GetUserTeams(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(teams)
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
func (h *TeamHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
    teamIDStr := chi.URLParam(r, "id")
    teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid team id", http.StatusBadRequest)
        return
    }

    inviterID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var req entity.InviteRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.service.InviteUser(teamID, inviterID, req.UserID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

