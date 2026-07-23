package controller

import (
	"github.com/gin-gonic/gin"
)

func RegisterTeamRoutes(r *gin.Engine, handler *TeamHandler) {
	api := r.Group("/api/v1")
	{
		api.POST("/teams", handler.CreateTeam)
		api.GET("/teams", handler.GetUserTeams)
		api.POST("/teams/:id/invite", handler.InviteUser)
	}
}