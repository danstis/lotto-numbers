package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danstis/lotto-numbers/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetLotteryNumbers(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		wantStatusCode  int
		wantLines       int
		wantNumPerLine  int
		wantSubset      []int
		wantErrContains string
	}{
		{
			name:            "Default parameters",
			query:           "",
			wantStatusCode:  http.StatusOK,
			wantLines:       5,
			wantNumPerLine:  6,
			wantErrContains: "",
		},
		{
			name:            "Valid numbers list",
			query:           "?numbersList=1,2,3,4,5,6,7,8,9,10",
			wantStatusCode:  http.StatusOK,
			wantLines:       5,
			wantNumPerLine:  6,
			wantSubset:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantErrContains: "",
		},
		{
			name:            "Invalid numbers list",
			query:           "?numbersList=invalid",
			wantStatusCode:  http.StatusBadRequest, // Expecting a bad request status when numbersList is invalid
			wantLines:       0,                     // When there's an error, no lines should be returned
			wantNumPerLine:  0,                     // When there's an error, numPerLine should be irrelevant
			wantSubset:      nil,                   // No subset should be expected on error
			wantErrContains: "Bad request",
		},
		{
			name:            "Valid lines parameter",
			query:           "?lines=3",
			wantStatusCode:  http.StatusOK,
			wantLines:       3,
			wantNumPerLine:  6,
			wantErrContains: "",
		},
		{
			name:            "Valid numPerLine parameter",
			query:           "?numPerLine=5",
			wantStatusCode:  http.StatusOK,
			wantLines:       5,
			wantNumPerLine:  5,
			wantErrContains: "",
		},
		{
			name:            "Non-positive lines parameter",
			query:           "?lines=0",
			wantStatusCode:  http.StatusBadRequest,
			wantLines:       0,
			wantNumPerLine:  0,
			wantSubset:      nil,
			wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
		},
		{
			name:            "Non-positive numPerLine parameter",
			query:           "?numPerLine=0",
			wantStatusCode:  http.StatusBadRequest,
			wantLines:       0,
			wantNumPerLine:  0,
			wantSubset:      nil,
			wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
		},
		{
			name:            "Not enough numbers for numPerLine",
			query:           "?numbersList=1,2,3&numPerLine=5",
			wantStatusCode:  http.StatusInternalServerError,
			wantLines:       0,
			wantNumPerLine:  0,
			wantSubset:      nil,
			wantErrContains: "ensure the 'numbersList' contains enough unique numbers",
		},
		{
			name:            "Zero lines requested",
			query:           "?lines=0",
			wantStatusCode:  http.StatusBadRequest,
			wantLines:       0,
			wantNumPerLine:  0,
			wantSubset:      nil,
			wantErrContains: "'lines' and 'numPerLine' must be positive numbers",
		},
		{
			name:            "Insufficient numbers for numPerLine",
			query:           "?numbersList=1,2,3&numPerLine=4",
			wantStatusCode:  http.StatusInternalServerError,
			wantLines:       0,
			wantNumPerLine:  0,
			wantSubset:      nil,
			wantErrContains: "ensure the 'numbersList' contains enough unique numbers",
		},
	}

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(VersionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedVersion := "0.0.0-development" // This should be the current version set in version.go
	if rr.Body.String() != expectedVersion {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedVersion)
	}
}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/lottery-numbers"+tc.query, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetLotteryNumbers)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatusCode, rr.Code, "handler returned wrong status code")
			// For the invalid numbers list, we expect a plain text error message, not a JSON response
			if tc.wantStatusCode != http.StatusOK {
				assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"), "content type is not text/plain; charset=utf-8")
				assert.Contains(t, rr.Body.String(), tc.wantErrContains, fmt.Sprintf("response body should contain '%s'", tc.wantErrContains))
			} else {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not application/json")

				var lotteryNumbers models.LotteryNumbers
				err = json.NewDecoder(rr.Body).Decode(&lotteryNumbers)
				assert.NoError(t, err, "response body should be a valid LotteryNumbers JSON")
				assert.Len(t, lotteryNumbers.Lines, tc.wantLines, "there should be the correct number of lines of lottery numbers")
				for _, line := range lotteryNumbers.Lines {
					assert.Len(t, line, tc.wantNumPerLine, "each line should contain the correct number of numbers")
					if tc.wantSubset != nil {
						assert.Subset(t, tc.wantSubset, line, "each line should only contain numbers from the provided numbersList")
					}
				}
			}
		})
	}
}
