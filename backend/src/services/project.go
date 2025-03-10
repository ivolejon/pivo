package services

import (
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

type (
	meta              = map[string]any
	AddDocumentParams struct {
		Content  string
		FileName string
	}
)

type ProjectService struct {
	clientID    uuid.UUID
	model       string
	vectorStore *repositories.VectorStore
}

func NewProjectService(clientID uuid.UUID) *ProjectService {
	return &ProjectService{
		vectorStore: nil,
		model:       "",
		clientID:    clientID,
	}
}

func (c *ProjectService) Init(LLMmodelName string) error {
	// TODO: Break out the selection of the model into a separate function
	supportedOllamModels := []string{"ollama:llama3.2"}
	if !slices.Contains(supportedOllamModels, LLMmodelName) {
		return errors.New("Model not supported")
	}

	c.model = LLMmodelName
	llm, err := ollama.New(ollama.WithModel(strings.Split(LLMmodelName, ":")[1]))
	if err != nil {
		return errors.New("Error creating LLM")
	}
	store, err := repositories.NewVectorStore(llm, c.clientID)
	if err != nil {
		return errors.New("Error creating VectorStore")
	}
	c.vectorStore = store
	return nil
}

func (c *ProjectService) AddDocument(params AddDocumentParams) error {
	if c.vectorStore == nil {
		return errors.New("ProjectService not initialized, call Init() first")
	}
	_, err := c.vectorStore.Provider.AddDocuments([]schema.Document{
		{PageContent: params.Content, Metadata: meta{"filename": params.FileName}},
	})
	if err != nil {
		return err
	}
	return nil
}
