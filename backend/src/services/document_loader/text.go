package document_loader

import (
	"bytes"
	"context"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type TextLoader struct{}

func (p *TextLoader) toDocuments(data []byte, spliter textsplitter.TextSplitter) ([]schema.Document, error) {
	reader := bytes.NewReader(data)
	TEXT := documentloaders.NewText(reader)
	return TEXT.LoadAndSplit(context.Background(), spliter)
}
