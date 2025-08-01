package knowledge_base_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/knowledge_base"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

var (
	clientID  = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")
	projectID = uuid.MustParse("8f2b7acc-6321-11f0-80c8-eb9676f528c1")
)

func TestKnowledgeBaseServiceNew(t *testing.T) {
	_, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInit_ollama_llama3_2(t *testing.T) {
	client, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
	err = client.Init("ollama-llama3.2")
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInit_ollama_gemma3_27b(t *testing.T) {
	client, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
	err = client.Init("ollama-gemma3:27b")
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInitWithFaultModel(t *testing.T) {
	svc, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
	err = svc.Init("ollama-llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestKnowledgeBaseServiceAddDocument_gemma3_27b(t *testing.T) {
	clientID := uuid.New()
	projectID := uuid.New()
	// Initialize the KnowledgeBaseService with a new clientID and projectID

	svc, errSvc := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, errSvc)
	_ = svc.Init("ollama-gemma3:27b")

	filename := "ivo.txt"
	projectId := uuid.New()

	params := knowledge_base.AddDocumentParams{
		Documents: []schema.Document{
			{
				PageContent: "The color of the house on the hill is blue.",
				Metadata:    map[string]any{"filename": filename},
			},
		},
		Filename:  filename,
		ProjectID: projectId,
		Title:     "New doc hello",
	}

	_, err := svc.AddDocuments(params)
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceAddDocument_llama3_2(t *testing.T) {
	svc, errSvc := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, errSvc)
	_ = svc.Init("ollama-llama3.2")

	filename := "ivo.txt"
	projectId := uuid.New()

	params := knowledge_base.AddDocumentParams{
		Documents: []schema.Document{
			{
				PageContent: "The color of the house on the hill is blue.",
				Metadata:    map[string]any{"filename": filename},
			},
		},
		Filename:  filename,
		ProjectID: projectId,
		Title:     "New doc hello",
	}

	_, err := svc.AddDocuments(params)
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceAddDocumentNoInit(t *testing.T) {
	svc, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
	_, err = svc.AddDocuments(knowledge_base.AddDocumentParams{})
	require.Equal(t, "KnowledgeBaseService not initialized, call Init() first", err.Error())
}

func TestKnowledgeBaseServiceQuery(t *testing.T) {
	svc, err := knowledge_base.NewKnowledgeBaseService(clientID, projectID)
	require.NoError(t, err)
	err = svc.Init("ollama-gemma3:27b")
	require.NoError(t, err)

	docs := []schema.Document{
		{
			PageContent: "The color of the buss is yellow.",
			Metadata:    map[string]any{"filename": "ivo.txt"},
		},
	}
	params := knowledge_base.AddDocumentParams{
		Documents: docs,
		Filename:  "ivo.txt",
		ProjectID: uuid.New(),
		Title:     "New doc hello",
	}

	_, err = svc.AddDocuments(params)
	require.NoError(t, err)
	res, err := svc.Query("What is a lion? And what color is the buss?")
	require.NoError(t, err)

	require.Contains(t, strings.ToLower(*res), "cat")
	require.Contains(t, strings.ToLower(*res), "yellow")
}
