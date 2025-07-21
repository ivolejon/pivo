package web

import (
	"log"

	"github.com/gin-gonic/gin"
	streamer_svc "github.com/ivolejon/pivo/services/streamer"
)

func SetupWebsocket(r *gin.Engine) {
	hub := streamer_svc.CreateStreamHub()

	go hub.Run() // Start the hub in a goroutine

	r.GET("/ws", func(c *gin.Context) {
		log.Println("WebSocket connection request received")
		streamer_svc.ServeWs(hub, c)
	})
}
