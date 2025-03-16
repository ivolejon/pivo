package knowledge_base_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/knowledge_base"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

var clientID = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")

func TestKnowledgeBaseServiceNew(t *testing.T) {
	_ = knowledge_base.NewKnowledgeBaseService(clientID)
}

func TestKnowledgeBaseServiceInit(t *testing.T) {
	client := knowledge_base.NewKnowledgeBaseService(clientID)
	err := client.Init("ollama:llama3.2")
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInitWithFaultModel(t *testing.T) {
	svc := knowledge_base.NewKnowledgeBaseService(clientID)
	err := svc.Init("ollama:llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestKnowledgeBaseServiceAddDocument(t *testing.T) {
	svc := knowledge_base.NewKnowledgeBaseService(clientID)
	_ = svc.Init("ollama:llama3.2")
	_, err := svc.AddDocuments([]schema.Document{
		{
			PageContent: "The color of the house on the hill is blue.",
			Metadata:    map[string]any{"filename": "ivo.txt"},
		},
	})
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceAddDocumentNoInit(t *testing.T) {
	svc := knowledge_base.NewKnowledgeBaseService(clientID)
	_, err := svc.AddDocuments([]schema.Document{})
	require.Equal(t, "KnowledgeBaseService not initialized, call Init() first", err.Error())
}

func TestKnowledgeBaseServiceQuery(t *testing.T) {
	svc := knowledge_base.NewKnowledgeBaseService(clientID)
	err := svc.Init("ollama:llama3.2")
	require.NoError(t, err)

	docs := []schema.Document{
		{
			PageContent: "The color of the buss is yellow.",
			Metadata:    map[string]any{"filename": "ivo.txt"},
		},
	}

	_, err = svc.AddDocuments(docs)
	require.NoError(t, err)
	res, err := svc.Query("Who is Donald Trump? And what color is the buss?")
	require.NoError(t, err)

	// require.Contains(t, strings.ToLower(*res), "#") // Test for markdown notation
	require.Contains(t, strings.ToLower(*res), "donald")
	require.Contains(t, strings.ToLower(*res), "trump")
	require.Contains(t, strings.ToLower(*res), "yellow")
}
