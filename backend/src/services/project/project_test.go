package services_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	projectSvc "github.com/ivolejon/pivo/services/project"
	"github.com/stretchr/testify/require"
)

var clientID = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")

func TestProjectServiceNew(t *testing.T) {
	_ = projectSvc.NewProjectService(clientID)
}

func TestProjectServiceInit(t *testing.T) {
	client := projectSvc.NewProjectService(clientID)
	err := client.Init("ollama:llama3.2")
	require.NoError(t, err)
}

func TestProjectServiceInitWithFaultModel(t *testing.T) {
	svc := projectSvc.NewProjectService(clientID)
	err := svc.Init("ollama:llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestProjectServiceAddDocument(t *testing.T) {
	svc := projectSvc.NewProjectService(clientID)
	_ = svc.Init("ollama:llama3.2")
	err := svc.AddDocument(projectSvc.AddDocumentParams{
		Content:  "The color of the house on the hill is blue.",
		FileName: "test.txt",
	})
	require.NoError(t, err)
}

func TestProjectServiceAddDocumentNoInit(t *testing.T) {
	svc := projectSvc.NewProjectService(clientID)
	err := svc.AddDocument(projectSvc.AddDocumentParams{})
	require.Equal(t, "ProjectService not initialized, call Init() first", err.Error())
}

func TestProjectServiceQuery(t *testing.T) {
	svc := projectSvc.NewProjectService(clientID)
	err := svc.Init("ollama:llama3.2")
	require.NoError(t, err)
	_ = svc.AddDocument(projectSvc.AddDocumentParams{
		Content:  "The color of the buss is yellow.",
		FileName: "test.txt",
	})
	res, err := svc.Query("Who is Donald Trump? And what color is the buss?")
	require.NoError(t, err)

	require.Contains(t, strings.ToLower(*res), "#") // Test for markdown notation
	require.Contains(t, strings.ToLower(*res), "donald")
	require.Contains(t, strings.ToLower(*res), "trump")
	require.Contains(t, strings.ToLower(*res), "yellow")
}
