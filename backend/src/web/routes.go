package web

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/document_loader"
	"github.com/ivolejon/pivo/services/knowledge_base"
)

func SetupDefaultRoutes(r *gin.Engine) {
	defaultGroup := r.Group("")

	defaultGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	defaultGroup.POST("/knowledge", handleAddDocumentToKnowledgeBase)
	defaultGroup.POST("/project/question", handleQuestionAboutDocument)
}

func handleAddDocumentToKnowledgeBase(c *gin.Context) {
	clientId := uuid.New()
	documentLoaderSvc := document_loader.NewDocumentLoaderService()
	knowledgeBaseSvc := knowledge_base.NewKnowledgeBaseService(clientId)
	err := knowledgeBaseSvc.Init("ollama:llama3.2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	docParams := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: filepath.Ext(header.Filename),
		ChunkSize:    500,
		Overlap:      50,
		Data:         data,
		MetaData:     map[string]any{"filename": header.Filename},
	}

	docs, err := documentLoaderSvc.LoadAsDocuments(docParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	docIDs, err := knowledgeBaseSvc.AddDocuments(docs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename, "inserted ids": docIDs,
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
