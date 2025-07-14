package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.Default()) // All origins allowed by default
	return engine
}
