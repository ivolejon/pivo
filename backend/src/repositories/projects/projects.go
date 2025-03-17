package projects

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ztrue/tracerr"
)

type ProjectsRepository struct {
	*Queries
	pool *pgxpool.Pool
}

func NewProjectsRepository() (*ProjectsRepository, error) {
	ctx := context.Background()
	db, err := db.ConnectAndGetPool(ctx)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &ProjectsRepository{
		Queries: New(),
		pool:    db.Pool,
	}, nil
}

func (r *ProjectsRepository) GetProjectsByClientId(clientId uuid.UUID) ([]Project, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)
	if errA != nil {
		return nil, tracerr.Wrap(errA)
	}
	defer conn.Release()
	return r.Queries.GetProjectsByClientId(ctx, conn, clientId)
}

func (r *ProjectsRepository) AddProject(args AddProjectParams) (Project, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)
	if errA != nil {
		return Project{}, tracerr.Wrap(errA)
	}
	defer conn.Release()
	return r.Queries.AddProject(ctx, conn, args)
}
