package repositories

import (
	"context"
	"errors"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaDB struct {
	store chroma.Store
	ctx   context.Context
}

type meta = map[string]any

var errAdd = errors.New("error adding document")

func NewChromaDB(embedder *embeddings.EmbedderImpl, collectionId uuid.UUID) (*ChromaDB, error) {
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
		store: store,
		ctx:   context.Background(),
	}, nil
}

func (p *ChromaDB) AddDocuments(documents []Document) error {
	_, err := p.store.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "Tokyo", Metadata: meta{"population": 9.7, "area": 622}, Score: 1},
	})
	if err != nil {
		return errAdd
	}
	return nil
}

func (p *ChromaDB) SimilaritySearch(search string) string {
	return "SimilaritySearch"
}

func (p *ChromaDB) RemoveCollection() bool {
	err := p.store.RemoveCollection()
	return err == nil
}
