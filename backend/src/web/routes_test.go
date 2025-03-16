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

func TestFileUploadRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupDefaultRoutes(router)

	testFilePath := "./test_data/pdf_file.pdf"

	// Create a multipart form with the test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(testFilePath)
	assert.NoError(t, err)
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(testFilePath))
	assert.NoError(t, err)
	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	// Create the POST request
	req, err := http.NewRequest("POST", "/knowledge", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "File uploaded successfully")
	assert.Contains(t, w.Body.String(), filepath.Base(testFilePath))
}
