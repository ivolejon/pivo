package services

import (
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	chain_store "github.com/ivolejon/pivo/chains"
	"github.com/ivolejon/pivo/repositories/vector_store"
	"github.com/ivolejon/pivo/services/ai"
	"github.com/tmc/langchaingo/llms"
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
	vectorStore *vector_store.VectorStore
	llm         llms.Model
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

	c.llm = llm

	store, err := vector_store.NewVectorStore("ChromaDb", llm, c.clientID)
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

func (c *ProjectService) Query(question string) (*string, error) {
	if c.vectorStore == nil {
		return nil, errors.New("ProjectService not initialized, call Init() first")
	}

	cs := chain_store.NewChainStore(c.vectorStore)
	aiSvc, err := ai.NewAiService()
	if err != nil {
		return nil, err
	}

	baseChain := cs.GetBaseDocumentChain(c.llm)
	formatAtProperDocumentChain := cs.GetFormatAsDocumentChain(c.llm)

	aiSvc.AddChain(baseChain)
	aiSvc.AddChain(formatAtProperDocumentChain)

	res, err := aiSvc.Run(question)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ProjectService) Refine(question string) (*string, error) {
	return nil, nil
}
