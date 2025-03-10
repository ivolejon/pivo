package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services"
	"github.com/stretchr/testify/require"
)

var clientID = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")

func TestProjectServiceNew(t *testing.T) {
	_ = services.NewProjectService(clientID)
}

func TestProjectServiceInit(t *testing.T) {
	client := services.NewProjectService(clientID)
	err := client.Init("ollama:llama3.2")
	require.NoError(t, err)
}

func TestProjectServiceInitWithFaultModel(t *testing.T) {
	svc := services.NewProjectService(clientID)
	err := svc.Init("ollama:llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestProjectServiceAddDocument(t *testing.T) {
	svc := services.NewProjectService(clientID)
	_ = svc.Init("ollama:llama3.2")
	err := svc.AddDocument(services.AddDocumentParams{
		Content:  "The color of the house on the hill is blue.",
		FileName: "test.txt",
	})
	require.NoError(t, err)
}

func TestProjectServiceAddDocumentNoInit(t *testing.T) {
	svc := services.NewProjectService(clientID)
	err := svc.AddDocument(services.AddDocumentParams{
		Content:  "This is a test",
		FileName: "test.txt",
	})
	require.Equal(t, "ProjectService not initialized, call Init() first", err.Error())
}

func TestProjectServiceQuery(t *testing.T) {
	svc := services.NewProjectService(clientID)
	err := svc.Init("ollama:llama3.2")
	require.NoError(t, err)
	svc.AddDocument(services.AddDocumentParams{
		Content:  "The color of the buss is yellow.",
		FileName: "test.txt",
	})
	res, err := svc.Query("Who is Donald Trump? And what color is the buss?")
	require.NoError(t, err)
	require.Equal(t, "hej", *res)
}
