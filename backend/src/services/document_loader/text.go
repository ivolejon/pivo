package document_loader

import (
	"bytes"
	"context"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/ztrue/tracerr"
)

type TextLoader struct{}

func (p *TextLoader) toDocuments(data []byte, spliter textsplitter.TextSplitter) ([]schema.Document, error) {
	reader := bytes.NewReader(data)
	TEXT := documentloaders.NewText(reader)
	documents, err := TEXT.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		return []schema.Document{}, tracerr.Wrap(err)
	}
	return documents, nil
}
