package repositories_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
)

func getOllama() (*ollama.LLM, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return nil, err
	}
	return llm, nil
}

func getOllamaEmbedder() (*embeddings.EmbedderImpl, error) {
	llm, err := getOllama()
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return embedder, nil
}

var testCollectionId = uuid.New()

func TestNewChromaDB(t *testing.T) {
	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	db, err := repositories.NewChromaDB(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
	}
	if db == nil {
		t.Errorf("ChromaDB is nil")
	}
}

func TestChromaDBAddDocuments(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}

	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	chroma, err := repositories.NewChromaDB(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}
	documents := []repositories.ClientDocument{
		{ID: uuid.New(), Content: "The color of the house is blue."},
		{ID: uuid.New(), Content: "The color of the car is red."},
		{ID: uuid.New(), Content: "The color of the desk is orange."},
	}
	errAdd := chroma.AddDocuments(documents)
	require.NoError(t, errAdd)

	result, err := chains.Run(
		context.TODO(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(chroma.Store, 1),
		),
		"What color is the desk?",
	)
	require.NoError(t, err)
	require.True(t, strings.Contains(result, "Orange."), "expected orange in result")
}

func TestRemoveCollection(t *testing.T) {
	llm, err := getOllama()
	if err != nil {
		t.Errorf("Error creating LLM: %v", err)
		return
	}
	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	store, err := repositories.NewChromaDB(llm, embedder, testCollectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}
	if !store.RemoveCollection() {
		t.Errorf("Error removing collection")
	}
}
