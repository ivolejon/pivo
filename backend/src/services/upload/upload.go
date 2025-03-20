package upload

import (
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/document_loader"
	"github.com/ivolejon/pivo/services/knowledge_base"
	"github.com/ztrue/tracerr"
)

type UploadService struct {
	clientID          uuid.UUID
	documentLoaderSvc *document_loader.DocumentLoaderService
	knowledgeBaseSvc  *knowledge_base.KnowledgeBaseService
}

func NewUploadService(clientID uuid.UUID) (*UploadService, error) {
	documentLoaderSvc, err := document_loader.NewDocumentLoaderService()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	knowledgeBaseSvc, err := knowledge_base.NewKnowledgeBaseService(clientID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	err = knowledgeBaseSvc.Init("ollama:llama3.2")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &UploadService{
		documentLoaderSvc: documentLoaderSvc,
		knowledgeBaseSvc:  knowledgeBaseSvc,
		clientID:          clientID,
	}, nil
}

type UploadFileParams struct {
	Data     []byte
	Filename string
}

func (svc *UploadService) Save(params UploadFileParams) (*[]uuid.UUID, error) {
	// TODO: LoadAsDocumentsParams could be provided instead ??
	docParams := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: filepath.Ext(params.Filename),
		ChunkSize:    500,
		Overlap:      50,
		Data:         params.Data,
		MetaData:     map[string]any{"filename": params.Filename},
	}

	// Splits and converts the data into documents
	docs, err := svc.documentLoaderSvc.LoadAsDocuments(docParams)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	addDocParams := knowledge_base.AddDocumentParams{
		Documents: docs,
		Filename:  params.Filename,
		ProjectID: svc.clientID,
		Title:     params.Filename,
	}

	// Adds the documents to the knowledge base(vector storage), and saves the document to database
	docIDs, err := svc.knowledgeBaseSvc.AddDocuments(addDocParams)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return docIDs, nil
}
