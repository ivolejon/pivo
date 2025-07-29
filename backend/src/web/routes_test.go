package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ivolejon/pivo/web"
	"github.com/stretchr/testify/require"
)

var clientID = uuid.MustParse("b15377e4-60f1-11f0-9ce3-834692c66f23")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity in this example.
		// In a production environment, you should restrict this to your frontend's domain.
		return true
	},
}

type formData struct {
	fields map[string]string
	file   *string
}

func prepareForm(data formData) (io.Reader, string) {
	// Create a multipart form with the test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	for key, value := range data.fields {
		if err := writer.WriteField(key, value); err != nil {
			panic("Error writing field to multipart form")
		}
	}

	if data.file == nil {
		return body, writer.FormDataContentType()
	}

	file, err := os.Open(*data.file)
	if err != nil {
		panic("Error reading file")
	}
	defer file.Close()

	part, _ := writer.CreateFormFile("file", filepath.Base(*data.file))
	if _, err := io.Copy(part, file); err != nil {
		panic("Error copying file content")
	}

	return body, writer.FormDataContentType()
}

func TestAddFileToKnowledgeBase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupProjectRoutes(router)

	filepath := "./test_data/pdf_file.pdf"

	data := formData{
		fields: map[string]string{
			"projectId": "8f2b7acc-6321-11f0-80c8-eb9676f528c1",
		},
		file: &filepath,
	}

	body, contentType := prepareForm(data)

	req, err := http.NewRequest("POST", "/project/add-document", body)
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
	web.SetupProjectRoutes(router)

	filepath := "./test_data/unsupported.zip"

	data := formData{
		fields: map[string]string{
			"projectId": "8f2b7acc-6321-11f0-80c8-eb9676f528c1",
		},
		file: &filepath,
	}

	body, contentType := prepareForm(data)

	req, err := http.NewRequest("POST", "/project/add-document", body)
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
	web.SetupProjectRoutes(router)

	question := map[string]string{
		"question":  "What is the capital of France?",
		"projectId": "8f2b7acc-6321-11f0-80c8-eb9676f528c1",
	}
	body, err := json.Marshal(question)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/project/question", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Paris")
}

func QuestionRequestFactory() *http.Request {
	question := map[string]string{
		"question":  "What is the capital of France?",
		"projectId": "8f2b7acc-6321-11f0-80c8-eb9676f528c1",
	}
	body, _ := json.Marshal(question)

	req, _ := http.NewRequest("POST", "/project/question", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestSendQuestionWithoutProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupProjectRoutes(router)

	question := map[string]string{
		"question": "What is the capital of France?",
	}
	body, err := json.Marshal(question)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/project/question", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCanConnectWebsocketStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	web.SetupProjectRoutes(router)
	web.SetupWebsocket(router)

	testServer := httptest.NewServer(router)
	defer testServer.Close() // Close the server when the test finishes

	wsURL := "ws" + testServer.URL[len("http"):] +
		fmt.Sprintf("/ws?client_id=%s", clientID.String())

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to dial WebSocket")
	defer conn.Close()

	// Assert the HTTP status code from the WebSocket handshake response
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "Expected StatusSwitchingProtocols")

	ch := make(chan string)

	go func() {
		defer close(ch)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return // Normal closure, exit the loop
				}
				require.NoError(t, err, "Failed to read message from WebSocket")
			}
			ch <- string(message)
		}
	}()

	go func() {
		req := QuestionRequestFactory()

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}()

	select {
	case msg := <-ch:
		require.NotEmpty(t, msg, "Expected to receive a message from the WebSocket")
		require.True(t, strings.Contains(msg, "Paris"),
			"Expected message to contain 'Paris', got: %s", msg)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message from WebSocket")
	}
}
