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
	"github.com/ivolejon/pivo/settings"
)

var clientID = settings.Environment().ClientID // Workaround for the client ID, assuming it's a UUID

func SetupDocumentRoutes(r *gin.Engine) {
	documentRoutes := r.Group("/document")
	documentRoutes.POST("/add-document", handleAddDocumentToKnowledgeBase)
	documentRoutes.GET("/list-documents", handleGetDocumentsByProjectId)
}

func handleAddDocumentToKnowledgeBase(c *gin.Context) {
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

func handleGetDocumentsByProjectId(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID"})
		return
	}

	projectSvc, err := projects_svc.NewProjectService(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	documents, err := projectSvc.ListDocumentsConnectedToProject(projectUUID)
	if err != nil {
		// Handle the error appropriately, no docs found is not an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if len(documents) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	c.JSON(http.StatusOK, documents)
}
