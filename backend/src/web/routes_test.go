package web_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ivolejon/pivo/web"
	"github.com/stretchr/testify/assert"
)

func TestSetupDefaultRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode) // Important for cleaner output in tests

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestAddFileToKnowledgeBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	body, contentType := prepareFile("./test_data/pdf_file.pdf")

	// Create the POST request
	req, err := http.NewRequest("POST", "/knowledge", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "File uploaded successfully")
	assert.Contains(t, w.Body.String(), "pdf_file.pdf")
}

func TestAddNonSupportedFileToKnowledgeBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	body, contentType := prepareFile("./test_data/unsupported.zip")

	// Create the POST request
	req, err := http.NewRequest("POST", "/knowledge", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "You can only add .pdf or .txt files")
}

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
	io.Copy(part, file)
	writer.Close()
	return body, writer.FormDataContentType()
}
