package web

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/document_loader"
	projects_svc "github.com/ivolejon/pivo/services/projects"
	"github.com/ivolejon/pivo/services/upload"
)

func SetupDefaultRoutes(r *gin.Engine) {
	projectGroup := r.Group("/project")
	projectGroup.POST("/create-project", handleCreateProject)
	projectGroup.POST("/add-document", handleAddDocumentToKnowledgeBase)
	projectGroup.POST("/question", handleQuestionAboutDocument)
	projectGroup.POST("/refine")
}

func handleAddDocumentToKnowledgeBase(c *gin.Context) {
	clientId := uuid.New()

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

	uploadSvc, err := upload.NewUploadService(clientId)
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

type QuestionPayload struct {
	Question string `json:"question" binding:"required"`
}

func handleQuestionAboutDocument(c *gin.Context) {
	var payload QuestionPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "question": payload})
}

func handleCreateProject(c *gin.Context) {
	clientId := uuid.New()

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
