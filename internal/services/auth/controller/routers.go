package controller

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, handler *AuthHandler) {
	api := r.Group("/api/v1")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}
}
