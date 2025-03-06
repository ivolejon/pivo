package main

import (
	"github.com/gin-gonic/gin"
)

func NewHTTPServer() *gin.Engine {
	engine := gin.Default()
	return engine
}
