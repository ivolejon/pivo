package web_test

import (
	"net/http"
	"net/http/httptest"
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

	// w = httptest.NewRecorder()
	// req, _ = http.NewRequest("POST", "/somepost", nil)

	// assert.Equal(t, http.StatusNotFound, w.Code)
}
