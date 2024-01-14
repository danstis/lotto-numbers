package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRoutes(t *testing.T) {
	r := mux.NewRouter()
	ctx := context.Background()
	SetupRoutes(ctx, r)

	tests := []struct {
		method       string
		description  string
		route        string
		expectedCode int
	}{
		{
			method:       "GET",
			description:  "Index route",
			route:        "/",
			expectedCode: http.StatusOK,
		},
		{
			method:       "GET",
			description:  "Lottery numbers route",
			route:        "/numbers",
			expectedCode: http.StatusOK,
		},
		{
			method:       "GET",
			description:  "Version route",
			route:        "/version",
			expectedCode: http.StatusOK,
		},
		{
			method:       "GET",
			description:  "Assets route",
			route:        "/assets/app.js",
			expectedCode: http.StatusOK,
		},
		{
			method:       "GET",
			description:  "Styles route",
			route:        "/assets/style.css",
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, test.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedCode {
			t.Errorf("handler returned wrong status code for route %s: got %v want %v, response: %s",
				test.route, status, test.expectedCode, rr.Body.String())
		}
	}
}
