package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danstis/lotto-numbers/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetLotteryNumbers(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		wantStatusCode int
		wantLines      int
		wantNumPerLine int
		wantSubset     []int
	}{
		{
			name:           "Default parameters",
			query:          "",
			wantStatusCode: http.StatusOK,
			wantLines:      5,
			wantNumPerLine: 6,
		},
		{
			name:           "Valid numbers list",
			query:          "?numbersList=1,2,3,4,5,6,7,8,9,10",
			wantStatusCode: http.StatusOK,
			wantLines:      5,
			wantNumPerLine: 6,
			wantSubset:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:           "Invalid numbers list",
			query:          "?numbersList=invalid",
			wantStatusCode: http.StatusBadRequest, // Expecting a bad request status when numbersList is invalid
			wantLines:      0, // When there's an error, no lines should be returned
			wantNumPerLine: 0, // When there's an error, numPerLine should be irrelevant
			wantSubset:     nil, // No subset should be expected on error
		},
		{
			name:           "Valid lines parameter",
			query:          "?lines=3",
			wantStatusCode: http.StatusOK,
			wantLines:      3,
			wantNumPerLine: 6,
		},
		{
			name:           "Valid numPerLine parameter",
			query:          "?numPerLine=5",
			wantStatusCode: http.StatusOK,
			wantLines:      5,
			wantNumPerLine: 5,
		},
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
			if tc.name == "Invalid numbers list" {
				assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"), "content type is not text/plain; charset=utf-8")
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