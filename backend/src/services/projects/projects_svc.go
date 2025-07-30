package projects_svc

import (
	"time"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/documents"
	"github.com/ivolejon/pivo/repositories/projects"
	"github.com/ztrue/tracerr"
)

type ProjectService struct {
	clientID     uuid.UUID
	projectRepo  *projects.ProjectsRepository
	documentRepo *documents.DocumentsRepository
}

func NewProjectService(clientID uuid.UUID) (*ProjectService, error) {
	projectRepo, err := projects.NewProjectsRepository()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	documentRepo, err := documents.NewDocumentsRepository()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &ProjectService{
		clientID:     clientID,
		projectRepo:  projectRepo,
		documentRepo: documentRepo,
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

func (svc *ProjectService) ListProjects() ([]projects.Project, error) {
	projects, err := svc.projectRepo.GetProjectsByClientId(svc.clientID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return projects, nil
}

func (svc *ProjectService) ListDocumentsConnectedToProject(projectId uuid.UUID) ([]documents.Document, error) {
	documents, err := svc.documentRepo.GetDocumentsByProjectId(projectId)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return documents, nil
}
