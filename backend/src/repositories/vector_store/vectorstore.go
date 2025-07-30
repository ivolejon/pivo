package vector_store

import (
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/embedders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/ztrue/tracerr"
)

type VectorStorageProvider interface {
	AddDocuments([]schema.Document) ([]string, error)
	RemoveDocument(string) error
	SimilaritySearch(string, int) ([]schema.Document, error)
	RemoveCollection() bool
	GetRetriver(int) vectorstores.Retriever
	Close()
}

type VectorStore struct {
	Provider     VectorStorageProvider
	CollectionId uuid.UUID
	embedder     *embeddings.EmbedderImpl
}

func NewVectorStore(storeType string, llm llms.Model, collectionId uuid.UUID) (*VectorStore, error) {
	embedder, errEm := embedders.GetEmbedderBgeLarge()
	if errEm != nil {
		return nil, tracerr.Wrap(errEm)
	}
	var err error
	var store VectorStorageProvider

	switch storeType {
	case "ChromaDb":
		store, err = NewChromaDB(llm, embedder, collectionId)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	case "PgVector":
		store, err = NewPgVector(llm, embedder, collectionId)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	default:
		return nil, tracerr.New("Vector store is not of a supported type")
	}

	return &VectorStore{
		Provider:     store,
		CollectionId: collectionId,
		embedder:     embedder,
	}, nil
}

func (v *VectorStore) AddDocuments(documents []schema.Document) ([]string, error) {
	return v.Provider.AddDocuments(documents)
}

func (v *VectorStore) SimilaritySearch(search string, numOfResults int) ([]schema.Document, error) {
	return v.Provider.SimilaritySearch(search, numOfResults)
}

func (v *VectorStore) RemoveCollection() bool {
	return v.Provider.RemoveCollection()
}

func (v *VectorStore) RemoveDocument(documentID string) error {
	return v.Provider.RemoveDocument(documentID)
}

func (v *VectorStore) Retriver(numOfDocs int) vectorstores.Retriever {
	return v.Provider.GetRetriver(numOfDocs)
}

func (v *VectorStore) Close() {
	v.Provider.Close()
}
