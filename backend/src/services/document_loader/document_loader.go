package document_loader

import (
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

func NewDocumentLoaderService() *DocumentLoaderService {
	return &DocumentLoaderService{}
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
		return tracerr.New("You can only add .pdf or .txt files")
	}
	if params.ChunkSize < 1 {
		return tracerr.New("ChunkSize are too low")
	}
	if params.Overlap < 1 {
		return tracerr.New("Overlap values are too low")
	}
	return nil
}
