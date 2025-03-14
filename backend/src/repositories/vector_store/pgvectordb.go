package vector_store

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"github.com/ztrue/tracerr"
)

type PgVectorDb struct {
	Store pgvector.Store
	ctx   context.Context
	pool  *pgxpool.Pool
}

func NewPgVectore(llm llms.Model, embedder *embeddings.EmbedderImpl, collectionId uuid.UUID) (*PgVectorDb, error) {
	ctx := context.Background()
	conn, err := db.ConnectAndGetPool(ctx)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	pool := conn.Pool
	store, err := pgvector.New(
		ctx,
		pgvector.WithConn(pool),
		pgvector.WithEmbedder(embedder),
		pgvector.WithPreDeleteCollection(true),
		pgvector.WithCollectionName(collectionId.String()),
	)
	if err != nil {
		return nil, err
	}
	return &PgVectorDb{
		Store: store,
		ctx:   ctx,
		pool:  pool,
	}, nil
}

func (p *PgVectorDb) GetRetriver(numOfDocs int) vectorstores.Retriever {
	return vectorstores.ToRetriever(p.Store, numOfDocs)
}

func (p *PgVectorDb) AddDocuments(documents []schema.Document) ([]string, error) {
	docIDs, err := p.Store.AddDocuments(context.Background(), documents)
	if err != nil {
		return []string{}, errAdd
	}
	return docIDs, nil
}

func (p *PgVectorDb) SimilaritySearch(search string, numOfResults int) ([]schema.Document, error) {
	result, err := p.Store.SimilaritySearch(context.Background(), search, numOfResults)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PgVectorDb) RemoveCollection() bool {
	c, err := p.pool.Acquire(p.ctx)
	if err != err {
		return false
	}
	tx, errTx := c.BeginTx(p.ctx, pgx.TxOptions{})
	if errTx != errTx {
		return false
	}
	err = p.Store.RemoveCollection(p.ctx, tx)
	return err == nil
}

func (p *PgVectorDb) RemoveDocument(id string) error {
	return nil
}
