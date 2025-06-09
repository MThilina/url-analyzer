package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-analyzer/internal/handler"
	"url-analyzer/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/analyze", handler.AnalyzeHandler)
	return r
}

func TestAnalyzeHandler_InvalidURL(t *testing.T) {
	router := setupRouter()

	reqBody := model.AnalyzeRequest{
		URL: "not-a-valid-url",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid URL")
}

func TestAnalyzeHandler_MissingURLField(t *testing.T) {
	router := setupRouter()

	reqBody := `{}` // missing "url" field
	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request")
}

func TestAnalyzeHandler_ValidURL(t *testing.T) {
	router := setupRouter()

	reqBody := model.AnalyzeRequest{
		URL: "https://example.com",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "title")
	assert.Contains(t, resp.Body.String(), "headings")
	assert.Contains(t, resp.Body.String(), "links")
}
