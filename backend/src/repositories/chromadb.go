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
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaDB struct {
	Store chroma.Store
	ctx   context.Context
}

type meta = map[string]any

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

func (p *ChromaDB) AddDocuments(documents []ClientDocument) error {
	_documents := []schema.Document{}
	for _, doc := range documents {
		_documents = append(_documents, schema.Document{
			PageContent: doc.Content,
			Metadata:    meta{},
			Score:       1,
		})
	}
	_, err := p.Store.AddDocuments(context.Background(), _documents)
	if err != nil {
		return errAdd
	}
	return nil
}

func (p *ChromaDB) SimilaritySearch(search string) string {
	return "SimilaritySearch"
}

func (p *ChromaDB) RemoveCollection() bool {
	err := p.Store.RemoveCollection()
	return err == nil
}
