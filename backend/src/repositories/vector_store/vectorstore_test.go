package vector_store_test

import (
	"testing"

	"github.com/ivolejon/pivo/repositories/vector_store"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func TestVectorStoreNew(t *testing.T) {
	llm, err := getOllama()
	require.NoError(t, err)

	vstore, err := vector_store.NewVectorStore("ChromaDb", llm, testCollectionId)
	require.NoError(t, err)
	defer vstore.Provider.Close()
}

func TestVectorStoreAddDocuments(t *testing.T) {
	llm, err := getOllama()
	require.NoError(t, err)

	store, err := vector_store.NewVectorStore("ChromaDb", llm, testCollectionId)
	require.NoError(t, err)
	defer store.Provider.Close()

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

	store, err := vector_store.NewVectorStore("ChromaDb", llm, testCollectionId)
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
