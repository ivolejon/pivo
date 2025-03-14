package settings

import (
	"fmt"
	"os"
	"sync"
)

type (
	webServerPort     = string
	databaseUrl       = string
	ChromaDatabaseUrl = string
)

type EnvironmentSettings struct {
	WebServerPort webServerPort
	DatabaseUrl   databaseUrl
	ChromaUrl     ChromaDatabaseUrl
}

var (
	instance *EnvironmentSettings
	once     sync.Once
)

func getChromaDbUrl() ChromaDatabaseUrl {
	chromaDbUrl := os.Getenv("CHROMA_URL")
	if chromaDbUrl == "" {
		panic("CHROMA_URL is missing")
	}
	return chromaDbUrl
}

func getWebServerPort() webServerPort {
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is missing")
	}
	return port
}

func getDbUrl() databaseUrl {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is missing")
	}
	return dbUrl
}

func Environment() *EnvironmentSettings {
	once.Do(func() {
		instance = &EnvironmentSettings{
			DatabaseUrl:   getDbUrl(),
			ChromaUrl:     getChromaDbUrl(),
			WebServerPort: fmt.Sprintf(":%s", getWebServerPort()),
		}
	})
	return instance
}
