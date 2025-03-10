package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// TODO: Should these be public?
type Provider interface {
	AddDocuments([]schema.Document) ([]string, error)
	RemoveDocument(string) error
	SimilaritySearch(string, int) ([]schema.Document, error)
	RemoveCollection() bool
	GetRetriver(int) vectorstores.Retriever
}

type VectorStore struct {
	Provider     Provider
	CollectionId uuid.UUID
	embedder     *embeddings.EmbedderImpl
}

func NewVectorStore(llm llms.Model, collectionId uuid.UUID) (*VectorStore, error) {
	// Here we can use muliple providers like ChromaDB, PgVector, etc. For now, we are using ChromaDB
	switch llm := llm.(type) {
	case *ollama.LLM:
		embedder, err := embeddings.NewEmbedder(llm)
		if err != nil {
			return nil, err
		}
		store, err := NewChromaDB(llm, embedder, collectionId)
		if err != nil {
			return nil, err
		}
		return &VectorStore{
			Provider:     store,
			CollectionId: collectionId,
			embedder:     embedder,
		}, nil
	default:
		return nil, errors.New("llm is not of a supported type")
	}
}

func (v *VectorStore) AddDocuments(documents []schema.Document) ([]string, error) {
	return v.Provider.AddDocuments(documents)
}

func (v *VectorStore) SimilaritySearch(search string, numOfResults int) ([]schema.Document, error) {
	return v.Provider.SimilaritySearch(search, numOfResults)
}

func (v *VectorStore) RemoveCollection(id uuid.UUID) bool {
	return v.Provider.RemoveCollection()
}

func (v *VectorStore) Retriver(numOfDocs int) vectorstores.Retriever {
	return v.Provider.GetRetriver(numOfDocs)
}
