package projects_svc

import (
	"time"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/projects"
	"github.com/ztrue/tracerr"
)

type ProjectService struct {
	clientID    uuid.UUID
	projectRepo *projects.ProjectsRepository
}

func NewProjectService(clientID uuid.UUID) (*ProjectService, error) {
	projectRepo, err := projects.NewProjectsRepository()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &ProjectService{
		clientID:    clientID,
		projectRepo: projectRepo,
	}, nil
}

func (svc *ProjectService) CreateNewProject(projectId uuid.UUID, title string) (*projects.Project, error) {
	newProjectParams := projects.AddProjectParams{
		ID:        projectId,
		ClientID:  svc.clientID,
		Title:     &title,
		CreatedAt: time.Now(),
	}

	project, err := svc.projectRepo.AddProject(newProjectParams)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return project, nil
}
