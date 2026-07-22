package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewAuthRouter(handler *AuthHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/api/v1/register", handler.Register)
	r.Post("/api/v1/login", handler.Login)

	return r
}

