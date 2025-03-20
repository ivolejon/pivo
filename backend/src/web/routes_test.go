package web_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ivolejon/pivo/web"
	"github.com/stretchr/testify/require"
)

func prepareFile(path string) (io.Reader, string) {
	// Create a multipart form with the test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(path)
	if err != nil {
		panic("Error reading file")
	}
	defer file.Close()

	part, _ := writer.CreateFormFile("file", filepath.Base(path))
	if _, err := io.Copy(part, file); err != nil {
		panic("Error copying file content")
	}
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestAddFileToKnowledgeBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	body, contentType := prepareFile("./test_data/pdf_file.pdf")

	req, err := http.NewRequest("POST", "/project/knowledge", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "File uploaded successfully")
	require.Contains(t, w.Body.String(), "pdf_file.pdf")
}

func TestAddNonSupportedFileToKnowledgeBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	body, contentType := prepareFile("./test_data/unsupported.zip")

	req, err := http.NewRequest("POST", "/project/knowledge", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "File type not supported.")
}

func TestSendQuestionToProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	question := map[string]string{"question": "What is the capital of France?"}
	body, err := json.Marshal(question)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/project/question", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "What is the capital of France?")
}
