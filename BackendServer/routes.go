package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRoutes(group *gin.RouterGroup) {
	health := group.Group("/health")
	{
		health.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
	}

	users := group.Group("/gamers")
	{
		users.POST("/points/credit", creditPoints)
	}
}

func creditPoints(ctx *gin.Context) {
	reqBody := UserGamePoints{}
	//time.Sleep(time.Millisecond * 10)
	if err := ctx.ShouldBind(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, reqBody)
}

type UserGamePoints struct {
	Game    string `json:"game"`
	Points  int    `json:"points"`
	GamerID string `json:"gamer_id"`
}
