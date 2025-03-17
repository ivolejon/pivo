package vector_store_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/embedders"
	"github.com/ivolejon/pivo/repositories/vector_store"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func TestNewPgVectorDb(t *testing.T) {
	embedder, err := embedders.GetEmbedderNomicEmbedTextModel()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	vstore, err := vector_store.NewPgVector(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating PgVector: %v", err)
	}
	if vstore == nil {
		t.Errorf("PgVector is nil")
	}
	defer vstore.Close()
}

func TestPgVectorDbAddDocuments(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	embedder, err := embedders.GetEmbedderNomicEmbedTextModel()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	vstore, err := vector_store.NewPgVector(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}
	defer vstore.Close()
	documents := []schema.Document{
		{PageContent: "The color of the house is blue."},
		{PageContent: "The color of the car is red."},
		{PageContent: "The color of the desk is orange."},
	}
	docIDs, errAdd := vstore.AddDocuments(documents)
	require.NoError(t, errAdd)
	require.Len(t, docIDs, 3)

	// TODO: Move this to a separate test
	result, err := chains.Run(
		context.TODO(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(vstore.Store, 1),
		),
		"What color is the house?",
	)
	require.NoError(t, err)
	require.True(t, strings.Contains(strings.ToLower(result), "blue"), "expected blue in result")
}

func TestPgVectorDbSimilaritySearch(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	embedder, err := embedders.GetEmbedderNomicEmbedTextModel()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	vstore, err := vector_store.NewPgVector(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}

	defer vstore.Close()

	docIDs, err := vstore.AddDocuments([]schema.Document{
		{PageContent: "tokyo", Metadata: map[string]any{
			"country": "japan", "id": uuid.New().String(),
			"filename": "tokyo.txt",
		}},
		{PageContent: "potato"},
	})

	require.NoError(t, err)
	require.Len(t, docIDs, 2)

	docs, err := vstore.SimilaritySearch("tokyo", 1)
	require.NoError(t, err)
	require.Len(t, docs, 1)
	require.Equal(t, "tokyo", docs[0].PageContent)
	country := docs[0].Metadata["country"]
	require.NoError(t, err)
	require.Equal(t, "japan", country)
}

func TestRemovePgVectorDbCollection(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}
	embedder, err := embedders.GetEmbedderNomicEmbedTextModel()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	vstore, err := vector_store.NewPgVector(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating PgVector: %v", err)
		return
	}
	defer vstore.Close()
	if !vstore.RemoveCollection() {
		t.Errorf("Error removing collection")
	}
}
