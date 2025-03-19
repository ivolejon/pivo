package documents

import (
	"context"

	"github.com/ivolejon/pivo/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ztrue/tracerr"
)

type DocumentsRepository struct {
	*Queries
	pool *pgxpool.Pool
}

func NewDocumentsRepository() (*DocumentsRepository, error) {
	ctx := context.Background()
	db, err := db.ConnectAndGetPool(ctx)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &DocumentsRepository{
		Queries: New(),
		pool:    db.Pool,
	}, nil
}

func (r *DocumentsRepository) AddDocument(args AddDocumentParams) (*Document, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)
	if errA != nil {
		return nil, tracerr.Wrap(errA)
	}
	defer conn.Release()
	document, err := r.Queries.AddDocument(ctx, conn, args)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &document, nil
}

// func (r *ProjectsRepository) GetProjectById(projectID uuid.UUID) (*Project, error) {
// 	ctx := context.Background()
// 	conn, errA := r.pool.Acquire(ctx)
// 	if errA != nil {
// 		return nil, tracerr.Wrap(errA)
// 	}
// 	defer conn.Release()
// 	project, err := r.Queries.GetProjectById(ctx, conn, projectID)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return nil, nil
// 		} else {
// 			return nil, err
// 		}
// 	}
// 	return &project, nil
// }
