package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	// Arrange
	called := false
	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	})
	handler := LoggingMiddleware(nextHandler)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	assert.True(t, called, "The next handler was not called")

	// Check if the log contains the expected values
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	handler.ServeHTTP(recorder, request)
	logOutput := buf.String()

	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/")
	assert.Contains(t, logOutput, "Request from")
}

func TestLoggingMiddlewareWithDifferentMethod(t *testing.T) {
	// Arrange
	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := LoggingMiddleware(nextHandler)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/test", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Check if the log contains the expected values
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	handler.ServeHTTP(recorder, request)
	logOutput := buf.String()

	assert.Contains(t, logOutput, "POST")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "Request from")
}
