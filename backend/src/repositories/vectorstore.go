package repositories

import (
	"github.com/google/uuid"
)

type Document struct {
	ID      uuid.UUID
	content string
}

type Provider interface {
	AddDocuments([]Document) error
	SimilaritySearch() string
	RemoveDocument(id uuid.UUID) bool
}

type VectorStore struct {
	Provider               Provider
	GetDocumentsByClientID func() []Document
}

func NewDocument(content string) Document {
	return Document{
		ID:      uuid.New(),
		content: content,
	}
}

func NewVectorStore() (*VectorStore, error) {
	// Here we can use muliple providers like ChromaDB, PgVector, etc. For now, we are using ChromaDB
	store, err := NewChromaDB()
	if err != nil {
		return nil, err
	}
	return &VectorStore{
		Provider: store,
	}, nil
}

func (v *VectorStore) AddDocument() error {
	var documents []Document
	return v.Provider.AddDocuments(documents)
}

func (v *VectorStore) SimilaritySearch() string {
	return v.Provider.SimilaritySearch()
}

func (v *VectorStore) RemoveDocument(id uuid.UUID) bool {
	return v.Provider.RemoveDocument(id)
}
