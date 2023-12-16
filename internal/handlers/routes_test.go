package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRoutes(t *testing.T) {
	r := mux.NewRouter()
	SetupRoutes(r)

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Index route",
			route:        "/",
			expectedCode: http.StatusOK,
		},
		{
			description:  "Lottery numbers route",
			route:        "/numbers",
			expectedCode: http.StatusOK,
		},
		{
			description:  "Assets route",
			route:        "/assets/app.js",
			expectedCode: http.StatusOK,
		},
		{
			description:  "Styles route",
			route:        "/assets/style.css",
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, test.expectedCode)
		}
	}
}
