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
	"github.com/tmc/langchaingo/prompts"
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

	systemTemplate := `
	You are an expert information retrieval and synthesis agent. Your primary task is to answer user queries accurately and comprehensively.

	**Process:**

	1.  **Prioritize Context:** Carefully analyze the information provided in the given context. Treat this context as your primary source of truth. If the context is insufficient, you can use the LLM to fill in the gaps.**
	2. **Be Concise:** Provide a clear and concise response to the user query. Avoid unnecessary information or verbosity.**
	3. **Be Accurate:** Never refer to the context or the LLM in your answer.**
	4. **This is super important: Your output and answers should be json format, in this format, [{"title": string, "content": string}] for the different sections.**

	**In essence, context first, LLM second.**
	This is the context: {{.input_documents}}\n\n
	And this is the question: {{.question}}`

	prompt := prompts.NewPromptTemplate(
		systemTemplate,
		[]string{"input_documents", "question"},
	)
	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(c.llm, prompt))
	retriever := c.vectorStore.Retriver(2)
	chain := chains.NewRetrievalQA(combineChain, retriever)

	res, err := chains.Run(context.Background(), chain, question)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
