package projects_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/projects"
	"github.com/stretchr/testify/require"
)

func TestCanInitlizeRepo(t *testing.T) {
	_, err := projects.NewProjectsRepository()
	require.NoError(t, err)
}

func TestGetProjectsByClientId(t *testing.T) {
	clientID := uuid.New()
	repo, err := projects.NewProjectsRepository()
	require.NoError(t, err)
	projects, err := repo.GetProjectsByClientId(clientID)
	require.NoError(t, err)
	require.Equal(t, 0, len(projects))
}

func TestGetProjectById(t *testing.T) {
	projectID := uuid.New()
	repo, err := projects.NewProjectsRepository()
	require.NoError(t, err)
	project, err := repo.GetProjectById(projectID)
	require.NoError(t, err)
	require.Nil(t, project)
}

func TestAddProject(t *testing.T) {
	clientID := uuid.New()
	projectID := uuid.New()
	title := "IvoPivo"
	now := time.Now()
	repo, err := projects.NewProjectsRepository()
	require.NoError(t, err)

	params := projects.AddProjectParams{
		ID:        projectID,
		ClientID:  clientID,
		Title:     &title,
		CreatedAt: now,
	}

	project, err := repo.AddProject(params)
	require.NoError(t, err)
	if !now.Equal(project.CreatedAt) {
		t.Errorf("Expected %v, got %v", now, project.CreatedAt)
	}
	require.Equal(t, title, *project.Title)
}
