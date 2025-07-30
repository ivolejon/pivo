package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/knowledge_base"
	projects_svc "github.com/ivolejon/pivo/services/projects"
	"github.com/ivolejon/pivo/settings"
)

var clientId = settings.Environment().ClientID // Assuming ClientID is a UUID string in the environment settings

// Example client ID, replace with actual logic to get client ID

type QuestionPayload struct {
	ProjectID uuid.UUID `json:"projectId" binding:"required"`
	Question  string    `json:"question" binding:"required"`
}

func SetupProjectRoutes(r *gin.Engine) {
	projectGroup := r.Group("/project")
	projectGroup.POST("/create-project", handleCreateProject)
	projectGroup.GET("/list-projects", handleListProjects)
	projectGroup.POST("/question", handleQuestionAboutDocument)
	projectGroup.POST("/refine")
}

func handleCreateProject(c *gin.Context) {
	var project struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectSvc, err := projects_svc.NewProjectService(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newProject, err := projectSvc.CreateNewProject(uuid.New(), project.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newProject)
}

func handleQuestionAboutDocument(c *gin.Context) {
	var payload QuestionPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	knowledgeBase, err := knowledge_base.NewKnowledgeBaseService(clientId, payload.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = knowledgeBase.Init("ollama-gemma3:27b")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	response, err := knowledgeBase.Query(payload.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "response": response})
}

func handleListProjects(c *gin.Context) {
	projectSvc, err := projects_svc.NewProjectService(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects, err := projectSvc.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}
