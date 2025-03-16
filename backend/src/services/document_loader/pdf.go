package document_loader

import (
	"bytes"
	"context"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/ztrue/tracerr"
)

type PdfLoader struct{}

func (p *PdfLoader) toDocuments(data []byte, spliter textsplitter.TextSplitter) ([]schema.Document, error) {
	reader := bytes.NewReader(data)
	PDF := documentloaders.NewPDF(reader, int64(len(data)))
	documents, err := PDF.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		return []schema.Document{}, tracerr.Wrap(err)
	}
	return documents, nil
}
