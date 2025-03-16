package vector_store

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/db"
	"github.com/jackc/pgx/v5"
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
	conn  *pgx.Conn
}

func NewPgVector(llm llms.Model, embedder *embeddings.EmbedderImpl, collectionId uuid.UUID) (*PgVectorDb, error) {
	ctx := context.Background()
	db, err := db.ConnectAndGetPool(ctx)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	conn, errA := db.Pool.Acquire(ctx)
	if errA != nil {
		return nil, tracerr.Wrap(err)
	}
	store, err := pgvector.New(
		ctx,
		pgvector.WithConn(conn),
		pgvector.WithEmbedder(embedder),
		pgvector.WithPreDeleteCollection(true),
		pgvector.WithCollectionName(collectionId.String()),
		pgvector.WithVectorDimensions(768),
		pgvector.WithHNSWIndex(16, 64, "vector_l2_ops"),
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &PgVectorDb{
		Store: store,
		ctx:   ctx,
		conn:  conn.Conn(),
	}, nil
}

func (p *PgVectorDb) GetRetriver(numOfDocs int) vectorstores.Retriever {
	return vectorstores.ToRetriever(p.Store, numOfDocs)
}

func (p *PgVectorDb) AddDocuments(documents []schema.Document) ([]string, error) {
	docIDs, err := p.Store.AddDocuments(context.Background(), documents)
	if err != nil {
		return []string{}, tracerr.Wrap(err)
	}
	return docIDs, tracerr.Wrap(err)
}

func (p *PgVectorDb) SimilaritySearch(search string, numOfResults int) ([]schema.Document, error) {
	result, err := p.Store.SimilaritySearch(context.Background(), search, numOfResults)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return result, nil
}

func (p *PgVectorDb) RemoveCollection() bool {
	tx, errTx := p.conn.BeginTx(p.ctx, pgx.TxOptions{})
	if errTx != errTx {
		return false
	}
	err := p.Store.RemoveCollection(p.ctx, tx)
	return err == nil
}

func (p *PgVectorDb) RemoveDocument(id string) error {
	return nil
}

func (p *PgVectorDb) Close() {
	p.conn.Close(p.ctx)
}
