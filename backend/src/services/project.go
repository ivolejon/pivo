package services

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
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

func (c *ProjectService) Query(question string) (*string, error) {
	if c.vectorStore == nil {
		return nil, errors.New("ProjectService not initialized, call Init() first")
	}
	// systemTemplate := `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

	// systemPrompt := prompts.NewSystemMessagePromptTemplate(systemTemplate, []string{})
	// combinedStuffQAChain := LoadStuffQA(llm)
	chain := chains.NewConversationalRetrievalQAFromLLM(c.llm,
		c.vectorStore.Provider.GetRetriver(1), memory.NewConversationBuffer())

	res, err := chains.Run(context.Background(), chain, question)
	if err != nil {
		return nil, err
	}
	return &res, nil
	// propmt := prompts.ChatPromptTemplate{
	// 	Template: `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.``

	// }
	// chains.New
}
