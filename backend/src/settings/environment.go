package settings

import (
	"fmt"
	"os"
	"sync"
)

type (
	webServerPort = string
	databaseUrl   = string
)

type EnvironmentSettings struct {
	WebServerPort webServerPort
	DatabaseUrl   databaseUrl
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

func getDbUrl() databaseUrl {
	port := os.Getenv("DATABASE_URL")
	if port == "" {
		panic("DATABASE_URL is missing")
	}
	return port
}

func Environment() *EnvironmentSettings {
	once.Do(func() {
		instance = &EnvironmentSettings{
			WebServerPort: fmt.Sprintf(":%s", getWebServerPort()),
			DatabaseUrl:   getDbUrl(),
		}
	})
	return instance
}
