package repositories

import (
	"context"
	"errors"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaDB struct {
	store chroma.Store
	ctx   context.Context
}

type meta = map[string]any

var errAdd = errors.New("error adding document")

func newEmbedder() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New()
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return embedder, nil
}

func NewChromaDB() (*ChromaDB, error) {
	embedder, err := newEmbedder()
	if err != nil {
		return nil, err
	}

	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithEmbedder(embedder),
		chroma.WithNameSpace(uuid.New().String()),
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

func (p *ChromaDB) SimilaritySearch() string {
	return "SimilaritySearch"
}

func (p *ChromaDB) RemoveDocument(id uuid.UUID) bool {
	return true
}
