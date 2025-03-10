package repositories

import (
	"context"
	"errors"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaDB struct {
	Store chroma.Store
	ctx   context.Context
}

var errAdd = errors.New("error adding document")

func NewChromaDB(llm llms.Model, embedder *embeddings.EmbedderImpl, collectionId uuid.UUID) (*ChromaDB, error) {
	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithEmbedder(embedder),
		chroma.WithNameSpace(collectionId.String()),
	)
	if err != nil {
		return nil, err
	}
	return &ChromaDB{
		Store: store,
		ctx:   context.Background(),
	}, nil
}

func (p *ChromaDB) GetRetriver(numOfDocs int) vectorstores.Retriever {
	return vectorstores.ToRetriever(p.Store, numOfDocs)
}

func (p *ChromaDB) AddDocuments(documents []schema.Document) ([]string, error) {
	docIDs, err := p.Store.AddDocuments(context.Background(), documents)
	if err != nil {
		return []string{}, errAdd
	}
	return docIDs, nil
}

func (p *ChromaDB) SimilaritySearch(search string, numOfResults int) ([]schema.Document, error) {
	result, err := p.Store.SimilaritySearch(context.Background(), search, numOfResults)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *ChromaDB) RemoveCollection() bool {
	err := p.Store.RemoveCollection()
	return err == nil
}

func (p *ChromaDB) RemoveDocument(id string) error {
	return nil
}
