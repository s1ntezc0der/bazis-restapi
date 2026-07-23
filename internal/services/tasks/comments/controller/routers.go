package controller

import (
	"github.com/gin-gonic/gin"

	"mkk_bazis/pkg/jwt"
	"mkk_bazis/pkg/middleware"
)

func RegisterCommentRoutes(r *gin.Engine, handler *CommentHandler, jwtConfig *jwt.JWTConfig) {
	api := r.Group("/api/v1")
	api.Use(middleware.GinAuthMiddleware(jwtConfig))
	{
		api.POST("/tasks/:id/comments", handler.AddComment)
		api.GET("/tasks/:id/comments", handler.GetComments)
	}
}