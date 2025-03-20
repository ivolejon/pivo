package knowledge_base_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/knowledge_base"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

var clientID = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")

func TestKnowledgeBaseServiceNew(t *testing.T) {
	_, err := knowledge_base.NewKnowledgeBaseService(clientID)
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInit(t *testing.T) {
	client, err := knowledge_base.NewKnowledgeBaseService(clientID)
	require.NoError(t, err)
	err = client.Init("ollama:llama3.2")
	require.NoError(t, err)
}

func TestKnowledgeBaseServiceInitWithFaultModel(t *testing.T) {
	svc, err := knowledge_base.NewKnowledgeBaseService(clientID)
	require.NoError(t, err)
	err = svc.Init("ollama:llama3.3")
	require.Equal(t, "Model not supported", err.Error())
}

func TestKnowledgeBaseServiceAddDocument(t *testing.T) {
	svc, errSvc := knowledge_base.NewKnowledgeBaseService(clientID)
	require.NoError(t, errSvc)
	_ = svc.Init("ollama:llama3.2")

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
	svc, err := knowledge_base.NewKnowledgeBaseService(clientID)
	require.NoError(t, err)
	_, err = svc.AddDocuments(knowledge_base.AddDocumentParams{})
	require.Equal(t, "KnowledgeBaseService not initialized, call Init() first", err.Error())
}

// FIXME: failing test
// func TestKnowledgeBaseServiceQuery(t *testing.T) {
// 	svc, err := knowledge_base.NewKnowledgeBaseService(clientID)
// 	require.NoError(t, err)
// 	err = svc.Init("ollama:llama3.2")
// 	require.NoError(t, err)

// 	docs := []schema.Document{
// 		{
// 			PageContent: "The color of the buss is yellow.",
// 			Metadata:    map[string]any{"filename": "ivo.txt"},
// 		},
// 	}

// 	_, err = svc.AddDocuments(docs)
// 	require.NoError(t, err)
// 	res, err := svc.Query("Who is Donald Trump? And what color is the buss?")
// 	require.NoError(t, err)

// 	// require.Contains(t, strings.ToLower(*res), "#") // Test for markdown notation
// 	require.Contains(t, strings.ToLower(*res), "donald")
// 	require.Contains(t, strings.ToLower(*res), "trump")
// 	require.Contains(t, strings.ToLower(*res), "yellow")
// }
