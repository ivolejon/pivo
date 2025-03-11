package document_loader_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	documentloaderSvc "github.com/ivolejon/pivo/services/document_loader"
	"github.com/stretchr/testify/require"
)

func loadPdfFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return bs
}

func TestDocumentLoader__New(t *testing.T) {
	loader := documentloaderSvc.PdfLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 500, 100)
	require.NoError(t, err)
	require.IsType(t, &documentloaderSvc.DocumentLoaderService{}, svc)
}

func TestDocumentLoader__LoadDocument(t *testing.T) {
	loader := documentloaderSvc.PdfLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 500, 100)
	require.NoError(t, err)
	data := loadPdfFile("./test_data/pdf_file.pdf")
	docs, err := svc.LoadAsDocuments(data)
	require.NoError(t, err)
	require.Equal(t, len(docs), 77)
}
