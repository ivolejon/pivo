package document_loader_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/ivolejon/pivo/services/document_loader"
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
	svc, err := documentloaderSvc.NewDocumentLoaderService()
	require.NoError(t, err)
	require.IsType(t, &documentloaderSvc.DocumentLoaderService{}, svc)
}

func TestDocumentLoader__LoadPDFDocument(t *testing.T) {
	svc, err := documentloaderSvc.NewDocumentLoaderService()
	require.NoError(t, err)
	data := loadFile("./test_data/pdf_file.pdf")
	params := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: "pdf",
		ChunkSize:    500,
		Overlap:      100,
		Data:         data,
		MetaData: map[string]any{
			"filename": "pdf_file.pdf",
		},
	}
	docs, err := svc.LoadAsDocuments(params)
	require.NoError(t, err)
	require.Equal(t, 77, len(docs))
}

func TestDocumentLoader__LoadTEXTDocument(t *testing.T) {
	svc, err := documentloaderSvc.NewDocumentLoaderService()
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	params := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: "text",
		ChunkSize:    300,
		Overlap:      50,
		Data:         data,
	}
	docs, err := svc.LoadAsDocuments(params)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
}

func TestDocumentLoader__ChunkSizeTooLow(t *testing.T) {
	svc, _ := documentloaderSvc.NewDocumentLoaderService()
	_, errL := svc.LoadAsDocuments(document_loader.LoadAsDocumentsParams{
		TypeOfLoader: "text",
		ChunkSize:    0,
		Overlap:      0,
		Data:         []byte{},
	})
	require.EqualError(t, errL, "ChunkSize are too low", errL)
}

func TestDocumentLoader__LoadDocumentWithFilename(t *testing.T) {
	svc, err := documentloaderSvc.NewDocumentLoaderService()
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	params := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: "text",
		ChunkSize:    300,
		Overlap:      50,
		Data:         data,
		MetaData: map[string]any{
			"filename": "text_file.txt",
		},
	}
	docs, err := svc.LoadAsDocuments(params)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
	filenameFromDoc, ok := docs[0].Metadata["filename"].(string)
	require.Equal(t, true, ok)
	require.Equal(t, "text_file.txt", filenameFromDoc)
}

func TestDocumentLoader__LoadDocumentWithNoFilename(t *testing.T) {
	svc, err := documentloaderSvc.NewDocumentLoaderService()
	require.NoError(t, err)
	data := loadFile("./test_data/text_file.txt")
	params := document_loader.LoadAsDocumentsParams{
		TypeOfLoader: "text",
		ChunkSize:    300,
		Overlap:      50,
		Data:         data,
	}
	docs, err := svc.LoadAsDocuments(params)
	require.NoError(t, err)
	require.Equal(t, 5, len(docs))
	filenameFromDoc, ok := docs[0].Metadata["filename"].(string)
	require.Equal(t, false, ok)
	require.Equal(t, "", filenameFromDoc)
}
