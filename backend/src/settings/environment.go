package settings

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type (
	webServerPort     = string
	databaseUrl       = string
	ChromaDatabaseUrl = string
	ClientID          = string
)

type EnvironmentSettings struct {
	WebServerPort webServerPort
	DatabaseUrl   databaseUrl
	ChromaUrl     ChromaDatabaseUrl
	ClientID      uuid.UUID
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

func getClientID() uuid.UUID {
	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		panic("CLIENT_ID is missing")
	}
	if parsedID, err := uuid.Parse(clientID); err != nil {
		panic(fmt.Sprintf("CLIENT_ID is not a valid UUID: %s", clientID))
	} else {
		return parsedID
	}
}

func Environment() *EnvironmentSettings {
	once.Do(func() {
		_ = godotenv.Load(".env.development") // Load environment variables from .env file

		instance = &EnvironmentSettings{
			DatabaseUrl:   getDbUrl(),
			ChromaUrl:     getChromaDbUrl(),
			WebServerPort: fmt.Sprintf(":%s", getWebServerPort()),
			ClientID:      getClientID(),
		}
	})
	return instance
}
