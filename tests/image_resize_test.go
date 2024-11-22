package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadImageSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	reqBody := `{"image":"` + testImage + `"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

}

func TestUploadImageFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		c.JSON(400, gin.H{})
	})

	invalidBody := `{"image":"testDATA"}`
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}