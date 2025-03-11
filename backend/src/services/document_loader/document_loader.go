package document_loader

import (
	"errors"
	"io"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type FileReader struct {
	Reader io.ReaderAt
	Size   int64
}

type DocumentLoader interface {
	toDocuments([]byte, textsplitter.TextSplitter) ([]schema.Document, error)
}

type DocumentLoaderService struct {
	loader   DocumentLoader
	splitter textsplitter.TextSplitter
}

func NewDocumentLoaderService(loader DocumentLoader, chunkSize int, overlap int) (*DocumentLoaderService, error) {
	if chunkSize < 1 || overlap < 1 {
		return nil, errors.New("ChunkSize or overlap values are too low")
	}
	// TODO: Think about if we could other multiple text-splitter
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(chunkSize),
		textsplitter.WithChunkOverlap(overlap),
	)
	return &DocumentLoaderService{
		loader:   loader,
		splitter: splitter,
	}, nil
}

func (svc *DocumentLoaderService) LoadAsDocuments(data []byte, filename *string) ([]schema.Document, error) {
	splitter := svc.splitter
	docs, err := svc.loader.toDocuments(data, splitter)
	if err != nil {
		return []schema.Document{}, err
	}
	if filename == nil {
		return docs, nil
	}
	for _, doc := range docs {
		doc.Metadata["filename"] = filename
	}
	return docs, nil
}
