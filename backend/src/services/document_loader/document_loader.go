package document_loader

import (
	"errors"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/ztrue/tracerr"
)

type DocumentLoader interface {
	toDocuments([]byte, textsplitter.TextSplitter) ([]schema.Document, error)
}

type DocumentLoaderService struct{}

type LoadAsDocumentsParams struct {
	TypeOfLoader string
	ChunkSize    int
	Overlap      int
	Data         []byte
	MetaData     map[string]any
}

var (
	ErrFileTypeNotSupported = errors.New("File type not supported.")
	ErrChunkSizeTooLow      = errors.New("ChunkSize are too low.")
	ErrOverlapTooLow        = errors.New("Overlap values are too low.")
)

func NewDocumentLoaderService() (*DocumentLoaderService, error) {
	return &DocumentLoaderService{}, nil
}

func (svc *DocumentLoaderService) LoadAsDocuments(params LoadAsDocumentsParams) ([]schema.Document, error) {
	err := validateLoadAsDocumentsParams(params)
	if err != nil {
		return []schema.Document{}, tracerr.Wrap(err)
	}
	var loader DocumentLoader
	switch params.TypeOfLoader {
	case "pdf", ".pdf":
		loader = &PdfLoader{}
	case "text", "txt", ".txt":
		loader = &TextLoader{}
	}

	// TODO: Think about if we could other multiple text-splitter
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(params.ChunkSize),
		textsplitter.WithChunkOverlap(params.Overlap),
	)
	docs, err := loader.toDocuments(params.Data, splitter)
	if err != nil {
		return []schema.Document{}, tracerr.Wrap(err)
	}
	if params.MetaData == nil {
		return docs, nil
	}
	// If there are any metedata, add it to the documents.
	for _, doc := range docs {
		for key, value := range params.MetaData {
			doc.Metadata[key] = value
		}
	}
	return docs, nil
}

func validateLoadAsDocumentsParams(params LoadAsDocumentsParams) error {
	allowedLoaders := map[string]bool{
		"pdf":  true,
		".pdf": true,
		"text": true,
		".txt": true,
	}

	if !allowedLoaders[params.TypeOfLoader] {
		return ErrFileTypeNotSupported
	}
	if params.ChunkSize < 1 {
		return ErrChunkSizeTooLow
	}
	if params.Overlap < 1 {
		return ErrOverlapTooLow
	}
	return nil
}
