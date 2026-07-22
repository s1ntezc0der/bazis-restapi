package controller

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

func NewTeamRouter(handler *TeamHandler) http.Handler {
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

    r.Post("/api/v1/teams", handler.CreateTeam)
    r.Get("/api/v1/teams", handler.GetUserTeams)
    r.Post("/api/v1/teams/{id}/invite", handler.InviteUser)

    return r
}

