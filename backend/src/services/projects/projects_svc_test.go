package projects_svc_test

import (
	"testing"

	"github.com/google/uuid"
	projects_svc "github.com/ivolejon/pivo/services/projects"
	"github.com/stretchr/testify/require"
)

var clientID = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")

func TestNewProjectService(t *testing.T) {
	_, err := projects_svc.NewProjectService(clientID)
	require.NoError(t, err)
}

func TestCreateNewProject(t *testing.T) {
	projectSvc, _ := projects_svc.NewProjectService(clientID)

	projectId := uuid.MustParse("067c673e-60c6-11f0-b83c-534d4bb70161")
	projectTitle := "New Project 1"

	newProject, err := projectSvc.CreateNewProject(projectId, projectTitle)
	require.NoError(t, err)
	require.Equal(t, newProject.ClientID, clientID)
	require.Equal(t, newProject.ID, projectId)
	require.Equal(t, newProject.Title, &projectTitle)
}
