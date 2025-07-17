package knowledge_base

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	chain_store "github.com/ivolejon/pivo/chains"
	"github.com/ivolejon/pivo/repositories/documents"
	"github.com/ivolejon/pivo/repositories/vector_store"
	"github.com/ivolejon/pivo/services/ai"
	"github.com/ivolejon/pivo/services/document_loader"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/ztrue/tracerr"
)

type (
	meta              = map[string]any
	AddDocumentParams struct {
		Documents []schema.Document
		Filename  string
		ProjectID uuid.UUID
		Title     string
	}
)

type KnowledgeBaseService struct {
	clientID          uuid.UUID
	projectID         uuid.UUID
	model             string
	vectorStore       *vector_store.VectorStore
	llm               llms.Model
	documentRepo      *documents.DocumentsRepository
	documentLoaderSvc *document_loader.DocumentLoaderService
}

func NewKnowledgeBaseService(clientID uuid.UUID, projectID uuid.UUID) (*KnowledgeBaseService, error) {
	documentRepo, err := documents.NewDocumentsRepository()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	documentLoaderSvc, err := document_loader.NewDocumentLoaderService()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &KnowledgeBaseService{
		vectorStore:       nil,
		model:             "",
		clientID:          clientID,
		projectID:         projectID,
		documentRepo:      documentRepo,
		documentLoaderSvc: documentLoaderSvc,
	}, nil
}

func (svc *KnowledgeBaseService) Init(LLMmodelName string) error {
	// TODO: Break out the selection of the model into a separate function
	supportedOllamModels := []string{"ollama:llama3.2"}
	if !slices.Contains(supportedOllamModels, LLMmodelName) {
		return errors.New("Model not supported")
	}

	svc.model = LLMmodelName
	llm, err := ollama.New(ollama.WithModel(strings.Split(LLMmodelName, ":")[1]))
	if err != nil {
		return tracerr.Wrap(err)
	}

	svc.llm = llm

	store, err := vector_store.NewVectorStore("ChromaDb", llm, svc.projectID)
	if err != nil {
		return tracerr.Wrap(err)
	}
	svc.vectorStore = store
	return nil
}

func (svc *KnowledgeBaseService) AddDocuments(params AddDocumentParams) (*[]uuid.UUID, error) {
	if svc.vectorStore == nil {
		return nil, errors.New("KnowledgeBaseService not initialized, call Init() first")
	}

	StringIds, err := svc.vectorStore.Provider.AddDocuments(params.Documents)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	// This should not happend.
	if len(StringIds) == 0 {
		return nil, tracerr.New("Missing IDs from vector store insertion.")
	}

	embeddingsIds := make([]uuid.UUID, len(StringIds))

	for i, v := range StringIds {
		embeddingsIds[i] = uuid.MustParse(v)
	}

	newDocParams := documents.AddDocumentParams{
		ID:            uuid.New(),
		EmbeddingsIds: embeddingsIds,
		Filename:      params.Filename,
		Title:         &params.Title,
		ProjectID:     params.ProjectID,
		CreatedAt:     time.Now(),
	}
	_, err = svc.documentRepo.AddDocument(newDocParams)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	// TODO: What should this return??
	return &embeddingsIds, nil
}

func (svc *KnowledgeBaseService) Query(question string) (*string, error) {
	if svc.vectorStore == nil {
		return nil, errors.New("KnowledgeBaseService not initialized, call Init() first")
	}

	cs := chain_store.NewChainStore(svc.vectorStore)
	aiSvc, err := ai.NewAiService()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	baseChain := cs.GetBaseDocumentChain(svc.llm)
	formatAtProperDocumentChain := cs.GetFormatAsDocumentChain(svc.llm)

	aiSvc.AddChain(baseChain)
	aiSvc.AddChain(formatAtProperDocumentChain)

	res, err := aiSvc.Run(question)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return res, nil
}

func (svc *KnowledgeBaseService) Refine(question string) (*string, error) {
	return nil, nil
}
