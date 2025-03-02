package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupStaticFilesRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router
	router := gin.Default()

	// Configure static file serving
	router.Static("/ui", "./ui/dist")
	router.StaticFile("/", "./ui/dist/index.html")

	return router
}

func TestStaticFilesServing(t *testing.T) {
	router := setupStaticFilesRouter()

	// Test cases for static files
	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "Index HTML",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html; charset=utf-8",
		},
		{
			name:           "CSS File",
			path:           "/ui/styles.css",
			expectedStatus: http.StatusOK,
			expectedType:   "text/css; charset=utf-8",
		},
		{
			name:           "JavaScript File",
			path:           "/ui/script.js",
			expectedStatus: http.StatusOK,
			expectedType:   "text/javascript; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, _ := http.NewRequest("GET", tc.path, nil)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			// Check the content type
			contentType := w.Header().Get("Content-Type")
			if contentType != tc.expectedType {
				t.Errorf("Expected content type %s, got %s", tc.expectedType, contentType)
			}
		})
	}
}
