package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupDefaultRoutes(r *gin.Engine) {
	defaultGroup := r.Group("")

	defaultGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
