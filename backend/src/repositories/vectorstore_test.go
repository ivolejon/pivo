package repositories_test

import (
	"testing"

	"github.com/ivolejon/pivo/repositories"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func TestVectorStoreNew(t *testing.T) {
	// embedder, err := getOllamaEmbedder()
	// if err != nil {
	// 	t.Errorf("Error creating embedder: %v", err)
	// 	return
	// }
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	_, err = repositories.NewVectorStore(llm, testCollectionId)
	if err != nil {
		t.Errorf("Error creating VectorStore: %v", err)
	}

	// db, err := repositories.NewChromaDB(llm, embedder, testCollectionId)
	// if err != nil {
	// 	t.Errorf("Error creating ChromaDB: %v", err)
	// }
	// if db == nil {
	// 	t.Errorf("ChromaDB is nil")
	// }
}

func TestVectorStoreAddDocuments(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	store, err := repositories.NewVectorStore(llm, testCollectionId)
	if err != nil {
		t.Errorf("Error creating VectorStore: %v", err)
	}

	docIDs, err := store.AddDocuments([]schema.Document{
		{
			PageContent: "Tokyo is the capital of Japan",
		},
	})
	require.NoError(t, err)
	require.Len(t, docIDs, 1)
}

func TestVectorStoreSimilaritySearch(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	store, err := repositories.NewVectorStore(llm, testCollectionId)
	if err != nil {
		t.Errorf("Error creating VectorStore: %v", err)
	}

	docIDs, err := store.AddDocuments([]schema.Document{
		{
			PageContent: "Tokyo is the capital of Japan",
		},
		{
			PageContent: "Stockholm is the capital of Sweden",
		},
	})
	require.NoError(t, err)
	require.Len(t, docIDs, 2)

	result, err := store.SimilaritySearch("Tokyo", 1)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "Tokyo is the capital of Japan", result[0].PageContent)
}
