package document_loader_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	documentloaderSvc "github.com/ivolejon/pivo/services/document_loader"
	"github.com/stretchr/testify/require"
)

func loadFile(path string) []byte {
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

func TestDocumentLoader__LoadPDFDocument(t *testing.T) {
	loader := documentloaderSvc.PdfLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 500, 100)
	require.NoError(t, err)
	data := loadFile("./test_data/pdf_file.pdf")
	docs, err := svc.LoadAsDocuments(data, nil)
	require.NoError(t, err)
	require.Equal(t, 77, len(docs))
}

func TestDocumentLoader__LoadTEXTDocument(t *testing.T) {
	loader := documentloaderSvc.TextLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 300, 50)
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	docs, err := svc.LoadAsDocuments(data, nil)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
}

func TestDocumentLoader__ChunkSizeTooLow(t *testing.T) {
	loader := documentloaderSvc.TextLoader{}
	_, err := documentloaderSvc.NewDocumentLoaderService(&loader, 0, 0)
	require.EqualError(t, err, "ChunkSize or overlap values are too low", err)
}

func TestDocumentLoader__LoadDocumentWithFilename(t *testing.T) {
	loader := documentloaderSvc.TextLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 300, 50)
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	var filename string = "pivo.txt"
	docs, err := svc.LoadAsDocuments(data, &filename)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
	filenameFromDoc, ok := docs[0].Metadata["filename"].(*string)
	require.Equal(t, ok, true)
	require.Equal(t, filename, *filenameFromDoc)
}

func TestDocumentLoader__LoadDocumentWithNoFilename(t *testing.T) {
	loader := documentloaderSvc.TextLoader{}
	svc, err := documentloaderSvc.NewDocumentLoaderService(&loader, 300, 50)
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	docs, err := svc.LoadAsDocuments(data, nil)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
	_, ok := docs[0].Metadata["filename"].(*string)
	require.Equal(t, ok, false)
}
