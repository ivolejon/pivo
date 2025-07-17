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

func (r *ProjectsRepository) GetProjectsByClientId(clientID uuid.UUID) ([]Project, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)
	if errA != nil {
		return nil, tracerr.Wrap(errA)
	}
	defer conn.Release()
	projects, err := r.Queries.GetProjectsByClientId(ctx, conn, clientID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return projects, nil
}

func (r *ProjectsRepository) AddProject(args AddProjectParams) (*Project, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)
	if errA != nil {
		return nil, tracerr.Wrap(errA)
	}
	defer conn.Release()
	project, err := r.Queries.AddProject(ctx, conn, args)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &project, nil
}

func (r *ProjectsRepository) GetProjectById(projectID uuid.UUID) (*Project, error) {
	ctx := context.Background()
	conn, errA := r.pool.Acquire(ctx)

	if errA != nil {
		return nil, tracerr.Wrap(errA)
	}
	defer conn.Release()
	project, err := r.Queries.GetProjectById(ctx, conn, projectID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &project, nil
}
