package repositories_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func getOllamaEmbedder() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return embedder, nil
}

var collectionId = uuid.New()

func TestNewChromaDB(t *testing.T) {
	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}

	db, err := repositories.NewChromaDB(embedder, collectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
	}
	if db == nil {
		t.Errorf("ChromaDB is nil")
	}
}

func TestChromaDBAddDocuments(t *testing.T) {
	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	store, err := repositories.NewChromaDB(embedder, collectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}
	errAdd := store.AddDocuments([]repositories.Document{})
	if errAdd != nil {
		t.Errorf("Error adding documents: %v", errAdd)
	}
}

func TestRemoveCollection(t *testing.T) {
	embedder, err := getOllamaEmbedder()
	if err != nil {
		t.Errorf("Error creating embedder: %v", err)
		return
	}
	store, err := repositories.NewChromaDB(embedder, collectionId)
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
		return
	}
	if !store.RemoveCollection() {
		t.Errorf("Error removing collection")
	}
}
