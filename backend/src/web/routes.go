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

	defaultGroup.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file") // "file" is the key of the form-data
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close() // Close the file when done

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": header.Filename})
	})
}
