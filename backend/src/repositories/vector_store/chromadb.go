package vector_store

import (
	"context"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/settings"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/ztrue/tracerr"
)

type ChromaDB struct {
	Store chroma.Store
	ctx   context.Context
}

func NewChromaDB(llm llms.Model, embedder *embeddings.EmbedderImpl, collectionId uuid.UUID) (*ChromaDB, error) {
	env := settings.Environment()
	store, err := chroma.New(
		chroma.WithChromaURL(env.ChromaUrl),
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
		return []string{}, tracerr.Wrap(err)
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

func (p *ChromaDB) Close() {
}
