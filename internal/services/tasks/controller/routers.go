package controller

import (
	"github.com/gin-gonic/gin"

	"mkk_bazis/pkg/jwt"
	"mkk_bazis/pkg/middleware"
)

func RegisterTaskRoutes(r *gin.Engine, handler *TaskHandler, jwtConfig *jwt.JWTConfig) {
	api := r.Group("/api/v1")
	api.Use(middleware.GinAuthMiddleware(jwtConfig))
	{
		api.POST("/tasks", handler.CreateTask)
		api.GET("/tasks", handler.GetTasks)
		api.PUT("/tasks/:id", handler.UpdateTask)
		api.GET("/tasks/:id/history", handler.GetHistory)
	}
}