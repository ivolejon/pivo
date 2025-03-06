package config

import (
	"fmt"
	"os"
	"sync"
)

type (
	webServerPort = string
	DbHost        = string
	DbPort        = int
	DbName        = string
	DbUsername    = string
	DbPassword    = string
)

type EnvironmentSettings struct {
	WebServerPort webServerPort
	DbHost        DbHost
	DbPort        DbPort
	DbName        DbName
	DbUsername    DbUsername
	DbPassword    DbPassword
}

var (
	instance *EnvironmentSettings
	once     sync.Once
)

func getWebServerPort() webServerPort {
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is missing")
	}
	return port
}

func Environment() *EnvironmentSettings {
	once.Do(func() {
		instance = &EnvironmentSettings{
			WebServerPort: fmt.Sprintf(":%s", getWebServerPort()),
			DbHost:        "",
			DbPort:        0,
			DbName:        "",
			DbUsername:    "",
			DbPassword:    "",
		}
	})
	return instance
}
