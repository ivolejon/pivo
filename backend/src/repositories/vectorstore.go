package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ClientDocument struct {
	ID      uuid.UUID
	Content string
}

type Provider interface {
	AddDocuments([]ClientDocument) error
	SimilaritySearch(string) string
	RemoveCollection() bool
}

type VectorStore struct {
	Provider Provider
	ClientId uuid.UUID
	embedder *embeddings.EmbedderImpl
	llm      llms.Model
}

// func NewDocument(content string) Document {
// 	return Document{
// 		ID:      uuid.New(),
// 		content: content,
// 	}
// }

type anyLLM = interface{}

func NewVectorStore(llm anyLLM, collectionId uuid.UUID) (*VectorStore, error) {
	// Here we can use muliple providers like ChromaDB, PgVector, etc. For now, we are using ChromaDB
	llm, ok := llm.(*ollama.LLM)
	if !ok {
		return nil, errors.New("llm is not of type *ollama.LLM")
	}
	ollama := llm.(*ollama.LLM)
	embedder, err := embeddings.NewEmbedder(ollama)
	if err != nil {
		return nil, err
	}

	store, err := NewChromaDB(ollama, embedder, collectionId)
	if err != nil {
		return nil, err
	}
	return &VectorStore{
		Provider: store,
		ClientId: collectionId,
		embedder: embedder,
		llm:      ollama,
	}, nil
}

func (v *VectorStore) AddDocuments(documents []ClientDocument) error {
	return v.Provider.AddDocuments(documents)
}

func (v *VectorStore) SimilaritySearch(search string) string {
	return v.Provider.SimilaritySearch(search)
}

func (v *VectorStore) RemoveCollection(id uuid.UUID) bool {
	return v.Provider.RemoveCollection()
}
