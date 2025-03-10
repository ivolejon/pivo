package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services"
	"github.com/stretchr/testify/require"
)

func TestProjectServiceNew(t *testing.T) {
	clientID := uuid.New()
	_ = services.NewProjectService(clientID)
}

func TestProjectServiceInit(t *testing.T) {
	clientID := uuid.New()
	client := services.NewProjectService(clientID)
	err := client.Init("ollama:llama3.2")
	require.NoError(t, err)
}

func TestProjectServiceInitWithFaultModel(t *testing.T) {
	clientID := uuid.New()
	svc := services.NewProjectService(clientID)
	err := svc.Init("ollama:llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestProjectServiceAddDocument(t *testing.T) {
	clientID := uuid.New()
	svc := services.NewProjectService(clientID)
	_ = svc.Init("ollama:llama3.2")
	err := svc.AddDocument(services.AddDocumentParams{
		Content:  "This is a test",
		FileName: "test.txt",
	})
	require.NoError(t, err)
}

func TestProjectServiceAddDocumentNoInit(t *testing.T) {
	clientID := uuid.New()
	svc := services.NewProjectService(clientID)
	err := svc.AddDocument(services.AddDocumentParams{
		Content:  "This is a test",
		FileName: "test.txt",
	})
	require.Equal(t, "ProjectService not initialized, call Init() first", err.Error())
}
