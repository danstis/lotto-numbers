package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSetupRoutes(t *testing.T) {
	r := mux.NewRouter()
	SetupRoutes(r)

	// Create a new HTTP request that will be tested against the router.
	req, err := http.NewRequest("GET", "/numbers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.Handler(r)

	// Serve the HTTP request to our router.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add more tests as needed to cover different routes and methods.
}
