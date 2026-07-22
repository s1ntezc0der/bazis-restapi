package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s1ntezc0der/bazis-restapi/pkg/jwt"
	"github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
)

func NewTaskRouter(handler *TaskHandler, jwtConfig *jwt.JWTConfig) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtConfig))

		r.Post("/api/v1/tasks", handler.CreateTask)
		r.Get("/api/v1/tasks", handler.GetTasks)
		r.Put("/api/v1/tasks/{id}", handler.UpdateTask)
		r.Get("/api/v1/tasks/{id}/history", handler.GetHistory)
	})

	return r
}

