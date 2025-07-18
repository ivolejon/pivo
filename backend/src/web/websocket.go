package web

import (
	"github.com/gin-gonic/gin"
	streamer_svc "github.com/ivolejon/pivo/services/streamer"
)

func SetupWebsocket(r *gin.Engine) {
	hub := streamer_svc.CreateStreamHub()

	go hub.Run() // Start the hub in a goroutine

	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		streamer_svc.ServeWs(hub, c)
	})
}
