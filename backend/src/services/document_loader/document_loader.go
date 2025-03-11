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
	// TODO: Think about if we could use multiple text-splitter
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(chunkSize),
		textsplitter.WithChunkOverlap(overlap),
	)
	return &DocumentLoaderService{
		loader:   loader,
		splitter: splitter,
	}, nil
}

func (svc *DocumentLoaderService) LoadAsDocuments(data []byte) ([]schema.Document, error) {
	splitter := svc.splitter
	return svc.loader.toDocuments(data, splitter)
}
