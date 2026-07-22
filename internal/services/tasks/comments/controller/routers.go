package controller

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/s1ntezc0der/bazis-restapi/pkg/jwt"
    "github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
)

func NewCommentRouter(handler *CommentHandler, jwtConfig *jwt.JWTConfig) http.Handler {
    r := chi.NewRouter()

    r.Group(func(r chi.Router) {
        r.Use(middleware.AuthMiddleware(jwtConfig))

        r.Post("/{id}/comments", handler.AddComment)
        r.Get("/{id}/comments", handler.GetComments)
    })

    return r
}
