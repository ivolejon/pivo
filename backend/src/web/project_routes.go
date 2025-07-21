package web

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/document_loader"
	"github.com/ivolejon/pivo/services/knowledge_base"
	projects_svc "github.com/ivolejon/pivo/services/projects"
	"github.com/ivolejon/pivo/services/upload"
)

var clientId = uuid.MustParse("b15377e4-60f1-11f0-9ce3-834692c66f23")

// Example client ID, replace with actual logic to get client ID

type QuestionPayload struct {
	ProjectID uuid.UUID `json:"projectId" binding:"required"`
	Question  string    `json:"question" binding:"required"`
}

func SetupProjectRoutes(r *gin.Engine) {
	projectGroup := r.Group("/project")
	projectGroup.POST("/create-project", handleCreateProject)
	projectGroup.POST("/add-document", handleAddDocumentToKnowledgeBase)
	projectGroup.POST("/question", handleQuestionAboutDocument)
	projectGroup.POST("/refine")
}

func handleAddDocumentToKnowledgeBase(c *gin.Context) {
	clientID := uuid.New()
	projectID := c.Request.FormValue("projectId")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	file, header, err := c.Request.FormFile("file") // "file" is the key of the form-data
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Prolems reading file"})
		return
	}

	uploadSvc, err := upload.NewUploadService(clientID, projectUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	IDs, err := uploadSvc.Save(upload.UploadFileParams{Data: data, Filename: header.Filename})
	if err != nil {
		if errors.Is(err, document_loader.ErrFileTypeNotSupported) ||
			errors.Is(err, document_loader.ErrChunkSizeTooLow) ||
			errors.Is(err, document_loader.ErrOverlapTooLow) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename, "inserted ids": IDs,
	})
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
	err = knowledgeBase.Init("ollama:llama3.2")
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
